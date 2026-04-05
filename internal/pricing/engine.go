// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

// Package pricing implements TerraKit's offline cost-estimation engine.
// All lookups are resolved entirely from a local pricing database — no network
// calls are ever made, which is the core privacy-first guarantee of TerraKit.
package pricing

import (
	"encoding/json"
	"fmt"
	"os"
)

// Estimate holds the result of a single resource cost lookup.
type Estimate struct {
	// HourlyUSD is the estimated on-demand cost per hour.
	HourlyUSD float64
	// MonthlyUSD is HourlyUSD × HoursPerMonth (730 h).
	MonthlyUSD float64
}

// Engine performs fully offline cost lookups against a PricingDB.
type Engine struct {
	db PricingDB
}

// NewEngine constructs a pricing engine.
//
// When pricingPath is non-empty, the JSON file at that path is loaded and
// merged on top of the built-in catalogue, allowing organisations to supply
// private negotiated rates or spot-price snapshots without sending data to any
// external service.
func NewEngine(pricingPath string) (*Engine, error) {
	db := builtinPricingDB()

	if pricingPath != "" {
		override, err := loadJSONDB(pricingPath)
		if err != nil {
			return nil, fmt.Errorf("loading external pricing db %q: %w", pricingPath, err)
		}
		// Merge: user-supplied entries win over built-ins.
		for resourceType, skus := range override {
			if db[resourceType] == nil {
				db[resourceType] = make(map[string]SKU)
			}
			for key, sku := range skus {
				db[resourceType][key] = sku
			}
		}
	}

	return &Engine{db: db}, nil
}

// Estimate returns the hourly and monthly USD costs for a given resource type
// and attribute map.
//
// For aws_instance, the "instance_type" attribute is used as the SKU key.
// For any other resource type, the "sku" attribute is used as a fallback.
func (e *Engine) Estimate(resourceType string, attrs map[string]string) (*Estimate, error) {
	skus, ok := e.db[resourceType]
	if !ok {
		return nil, fmt.Errorf("no pricing data found for resource type %q", resourceType)
	}

	key := skuKey(resourceType, attrs)
	if key == "" {
		return nil, fmt.Errorf("cannot determine SKU key for resource type %q — ensure the required attribute (e.g. instance_type) is set", resourceType)
	}

	sku, ok := skus[key]
	if !ok {
		return nil, fmt.Errorf("no SKU %q in pricing data for resource type %q", key, resourceType)
	}

	return &Estimate{
		HourlyUSD:  sku.HourlyUSD,
		MonthlyUSD: sku.HourlyUSD * HoursPerMonth,
	}, nil
}

// skuKey derives the canonical pricing-table key from the resource's attributes.
// Each resource type may use a different attribute as its discriminator.
func skuKey(resourceType string, attrs map[string]string) string {
	switch resourceType {
	// ── AWS: keyed by instance_type ───────────────────────────────────────
	case "aws_instance":
		if v, ok := attrs["instance_type"]; ok {
			return v
		}

	// ── AWS: keyed by instance_class ──────────────────────────────────────
	case "aws_db_instance", "aws_rds_cluster":
		if v, ok := attrs["instance_class"]; ok {
			return v
		}

	// ── AWS: EKS node group — instance_types is a list ────────────────────
	case "aws_eks_node_group":
		if v, ok := attrs["instance_types"]; ok {
			return v
		}
		if v, ok := attrs["instance_type"]; ok {
			return v
		}

	// ── AWS: ElastiCache keyed by node_type ───────────────────────────────
	case "aws_elasticache_cluster", "aws_elasticache_replication_group":
		if v, ok := attrs["node_type"]; ok {
			return v
		}

	// ── AWS: Redshift keyed by node_type ──────────────────────────────────
	case "aws_redshift_cluster":
		if v, ok := attrs["node_type"]; ok {
			return v
		}

	// ── AWS: OpenSearch / Elasticsearch keyed by instance_type ────────────
	case "aws_opensearch_domain", "aws_elasticsearch_domain":
		if v, ok := attrs["instance_type"]; ok {
			return v
		}

	// ── AWS: MSK keyed by broker instance type ────────────────────────────
	case "aws_msk_cluster":
		if v, ok := attrs["instance_type"]; ok {
			return v
		}

	// ── AWS: DynamoDB keyed by billing_mode ────────────────────────────────
	case "aws_dynamodb_table":
		if v, ok := attrs["billing_mode"]; ok {
			return v
		}
		return "PAY_PER_REQUEST"

	// ── AWS: Lambda keyed by memory_size ──────────────────────────────────
	case "aws_lambda_function":
		if v, ok := attrs["memory_size"]; ok {
			return v
		}

	// ── AWS: ECS keyed by launch_type ─────────────────────────────────────
	case "aws_ecs_service":
		if v, ok := attrs["launch_type"]; ok {
			return v
		}
		return "FARGATE"

	// ── AWS: EBS keyed by type ────────────────────────────────────────────
	case "aws_ebs_volume":
		if v, ok := attrs["type"]; ok {
			return v
		}
		return "gp3"

	// ── AWS: Kinesis keyed by stream_mode ─────────────────────────────────
	case "aws_kinesis_stream":
		if v, ok := attrs["stream_mode_details"]; ok {
			return v
		}
		return "ON_DEMAND"

	// ── AWS: Flat-rate resources ──────────────────────────────────────────
	case "aws_eks_cluster", "aws_eks_fargate_profile",
		"aws_nat_gateway", "aws_lb", "aws_alb",
		"aws_eip", "aws_vpc_endpoint",
		"aws_cloudfront_distribution", "aws_s3_bucket",
		"aws_sqs_queue", "aws_sns_topic",
		"aws_ecr_repository",
		"aws_ec2_transit_gateway", "aws_ec2_transit_gateway_vpc_attachment":
		return "standard"

	// ── Azure: keyed by size / vm_size ────────────────────────────────────
	case "azurerm_linux_virtual_machine", "azurerm_windows_virtual_machine":
		if v, ok := attrs["size"]; ok {
			return v
		}
	case "azurerm_virtual_machine":
		if v, ok := attrs["vm_size"]; ok {
			return v
		}

	// ── Azure: AKS keyed by sku_tier ──────────────────────────────────────
	case "azurerm_kubernetes_cluster":
		if v, ok := attrs["sku_tier"]; ok {
			return v
		}
		return "Free"

	// ── Azure: AKS node pool keyed by vm_size ─────────────────────────────
	case "azurerm_kubernetes_cluster_node_pool":
		if v, ok := attrs["vm_size"]; ok {
			return v
		}

	// ── Azure: SQL Database keyed by sku_name ─────────────────────────────
	case "azurerm_mssql_database":
		if v, ok := attrs["sku_name"]; ok {
			return v
		}

	// ── Azure: PostgreSQL/MySQL Flexible keyed by sku_name ────────────────
	case "azurerm_postgresql_flexible_server", "azurerm_mysql_flexible_server":
		if v, ok := attrs["sku_name"]; ok {
			return v
		}

	// ── Azure: Redis keyed by capacity (SKU family+capacity) ──────────────
	case "azurerm_redis_cache":
		if v, ok := attrs["sku_name"]; ok {
			if c, ok2 := attrs["capacity"]; ok2 {
				// e.g. "C0", "P1"
				return string(v[0]) + c
			}
		}

	// ── Azure: Managed Disk keyed by storage_account_type ─────────────────
	case "azurerm_managed_disk":
		if v, ok := attrs["storage_account_type"]; ok {
			return v
		}

	// ── Azure: App Service Plan keyed by sku_name ─────────────────────────
	case "azurerm_service_plan":
		if v, ok := attrs["sku_name"]; ok {
			return v
		}

	// ── Azure: Application Gateway keyed by sku name ──────────────────────
	case "azurerm_application_gateway":
		if v, ok := attrs["sku"]; ok {
			return v
		}

	// ── Azure: Container Registry keyed by sku ────────────────────────────
	case "azurerm_container_registry":
		if v, ok := attrs["sku"]; ok {
			return v
		}

	// ── Azure: VPN / ExpressRoute Gateway keyed by sku ────────────────────
	case "azurerm_virtual_network_gateway":
		if v, ok := attrs["sku"]; ok {
			return v
		}

	// ── Azure: Service Bus keyed by sku ───────────────────────────────────
	case "azurerm_servicebus_namespace":
		if v, ok := attrs["sku"]; ok {
			return v
		}

	// ── Azure: Flat-rate resources ────────────────────────────────────────
	case "azurerm_cosmosdb_account", "azurerm_public_ip",
		"azurerm_lb", "azurerm_nat_gateway", "azurerm_firewall",
		"azurerm_storage_account":
		return "standard"

	// ── GCP: Compute keyed by machine_type ────────────────────────────────
	case "google_compute_instance":
		if v, ok := attrs["machine_type"]; ok {
			return v
		}

	// ── GCP: GKE cluster keyed by mode ────────────────────────────────────
	case "google_container_cluster":
		// Check for autopilot config; default to standard
		if v, ok := attrs["enable_autopilot"]; ok && v == "true" {
			return "autopilot"
		}
		return "standard"

	// ── GCP: GKE node pool keyed by machine_type ──────────────────────────
	case "google_container_node_pool":
		if v, ok := attrs["machine_type"]; ok {
			return v
		}

	// ── GCP: Cloud SQL keyed by tier ──────────────────────────────────────
	case "google_sql_database_instance":
		if v, ok := attrs["tier"]; ok {
			return v
		}

	// ── GCP: Persistent Disk keyed by type ────────────────────────────────
	case "google_compute_disk":
		if v, ok := attrs["type"]; ok {
			return v
		}
		return "pd-balanced"

	// ── GCP: Cloud Functions keyed by available_memory_mb ──────────────────
	case "google_cloudfunctions_function", "google_cloudfunctions2_function":
		if v, ok := attrs["available_memory_mb"]; ok {
			return v
		}
		if v, ok := attrs["available_memory"]; ok {
			return v
		}

	// ── GCP: Flat-rate resources ──────────────────────────────────────────
	case "google_spanner_instance", "google_redis_instance",
		"google_compute_router_nat",
		"google_compute_forwarding_rule", "google_compute_global_forwarding_rule",
		"google_compute_address", "google_compute_global_address",
		"google_storage_bucket", "google_bigquery_dataset",
		"google_cloud_run_service", "google_cloud_run_v2_service",
		"google_pubsub_topic",
		"google_compute_vpn_gateway", "google_compute_ha_vpn_gateway",
		"google_artifact_registry_repository",
		"google_compute_security_policy":
		return "standard"
	}
	// Generic fallback: allow callers to pass an explicit "sku" attribute.
	return attrs["sku"]
}

// loadJSONDB reads and parses a user-supplied JSON pricing database file.
func loadJSONDB(path string) (PricingDB, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading pricing db file: %w", err)
	}
	var db PricingDB
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("parsing JSON pricing db: %w", err)
	}
	return db, nil
}
