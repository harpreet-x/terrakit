// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

package pricing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// ── NewEngine ────────────────────────────────────────────────────────────────

func TestNewEngine_BuiltinOnly(t *testing.T) {
	e, err := NewEngine("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestNewEngine_WithOverrideFile(t *testing.T) {
	db := PricingDB{
		"aws_instance": {
			"custom.xlarge": {HourlyUSD: 0.9999, Description: "custom"},
		},
	}
	path := writeTempPricingDB(t, db)
	e, err := NewEngine(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	est, err := e.Estimate("aws_instance", map[string]string{"instance_type": "custom.xlarge"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if est.HourlyUSD != 0.9999 {
		t.Errorf("hourly: got %v, want 0.9999", est.HourlyUSD)
	}
}

func TestNewEngine_OverrideWinsOverBuiltin(t *testing.T) {
	db := PricingDB{
		"aws_instance": {
			"t3.micro": {HourlyUSD: 1.0, Description: "overridden"},
		},
	}
	path := writeTempPricingDB(t, db)
	e, err := NewEngine(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	est, err := e.Estimate("aws_instance", map[string]string{"instance_type": "t3.micro"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if est.HourlyUSD != 1.0 {
		t.Errorf("override should win: got %v, want 1.0", est.HourlyUSD)
	}
}

func TestNewEngine_MissingFile(t *testing.T) {
	_, err := NewEngine("/nonexistent/pricing.json")
	if err == nil {
		t.Error("expected error for missing pricing file, got nil")
	}
}

func TestNewEngine_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte(`{ not json }`), 0o644)
	_, err := NewEngine(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// ── Estimate ─────────────────────────────────────────────────────────────────

func TestEstimate_MonthlyIsHourlyTimes730(t *testing.T) {
	e, _ := NewEngine("")
	est, err := e.Estimate("aws_instance", map[string]string{"instance_type": "t3.micro"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := est.HourlyUSD * HoursPerMonth
	if est.MonthlyUSD != want {
		t.Errorf("monthly: got %v, want %v", est.MonthlyUSD, want)
	}
}

func TestEstimate_UnknownResourceType(t *testing.T) {
	e, _ := NewEngine("")
	_, err := e.Estimate("aws_made_up_resource", map[string]string{})
	if err == nil {
		t.Error("expected error for unknown resource type, got nil")
	}
}

func TestEstimate_UnknownSKU(t *testing.T) {
	e, _ := NewEngine("")
	_, err := e.Estimate("aws_instance", map[string]string{"instance_type": "x9.does-not-exist"})
	if err == nil {
		t.Error("expected error for unknown SKU, got nil")
	}
}

func TestEstimate_MissingRequiredAttr(t *testing.T) {
	e, _ := NewEngine("")
	_, err := e.Estimate("aws_instance", map[string]string{})
	if err == nil {
		t.Error("expected error when required attribute is missing, got nil")
	}
}

// ── skuKey ───────────────────────────────────────────────────────────────────

func TestSkuKey_AWS(t *testing.T) {
	tests := []struct {
		rtype string
		attrs map[string]string
		want  string
	}{
		{"aws_instance", map[string]string{"instance_type": "t3.large"}, "t3.large"},
		{"aws_db_instance", map[string]string{"instance_class": "db.t3.micro"}, "db.t3.micro"},
		{"aws_eks_node_group", map[string]string{"instance_types": "m5.xlarge"}, "m5.xlarge"},
		{"aws_elasticache_cluster", map[string]string{"node_type": "cache.t3.micro"}, "cache.t3.micro"},
		{"aws_dynamodb_table", map[string]string{"billing_mode": "PROVISIONED"}, "PROVISIONED"},
		{"aws_dynamodb_table", map[string]string{}, "PAY_PER_REQUEST"},
		{"aws_lambda_function", map[string]string{"memory_size": "512"}, "512"},
		{"aws_ecs_service", map[string]string{}, "FARGATE"},
		{"aws_ecs_service", map[string]string{"launch_type": "EC2"}, "EC2"},
		{"aws_ebs_volume", map[string]string{}, "gp3"},
		{"aws_ebs_volume", map[string]string{"type": "io1"}, "io1"},
		{"aws_kinesis_stream", map[string]string{}, "ON_DEMAND"},
		{"aws_eks_cluster", map[string]string{}, "standard"},
		{"aws_nat_gateway", map[string]string{}, "standard"},
	}
	for _, tc := range tests {
		got := skuKey(tc.rtype, tc.attrs)
		if got != tc.want {
			t.Errorf("skuKey(%q, %v) = %q, want %q", tc.rtype, tc.attrs, got, tc.want)
		}
	}
}

func TestSkuKey_Azure(t *testing.T) {
	tests := []struct {
		rtype string
		attrs map[string]string
		want  string
	}{
		{"azurerm_linux_virtual_machine", map[string]string{"size": "Standard_D2s_v3"}, "Standard_D2s_v3"},
		{"azurerm_virtual_machine", map[string]string{"vm_size": "Standard_B2s"}, "Standard_B2s"},
		{"azurerm_kubernetes_cluster", map[string]string{}, "Free"},
		{"azurerm_kubernetes_cluster", map[string]string{"sku_tier": "Standard"}, "Standard"},
		{"azurerm_redis_cache", map[string]string{"sku_name": "C", "capacity": "0"}, "C0"},
		{"azurerm_managed_disk", map[string]string{"storage_account_type": "Premium_LRS"}, "Premium_LRS"},
		{"azurerm_cosmosdb_account", map[string]string{}, "standard"},
	}
	for _, tc := range tests {
		got := skuKey(tc.rtype, tc.attrs)
		if got != tc.want {
			t.Errorf("skuKey(%q, %v) = %q, want %q", tc.rtype, tc.attrs, got, tc.want)
		}
	}
}

func TestSkuKey_GCP(t *testing.T) {
	tests := []struct {
		rtype string
		attrs map[string]string
		want  string
	}{
		{"google_compute_instance", map[string]string{"machine_type": "n2-standard-2"}, "n2-standard-2"},
		{"google_container_cluster", map[string]string{}, "standard"},
		{"google_container_cluster", map[string]string{"enable_autopilot": "true"}, "autopilot"},
		{"google_sql_database_instance", map[string]string{"tier": "db-f1-micro"}, "db-f1-micro"},
		{"google_compute_disk", map[string]string{}, "pd-balanced"},
		{"google_compute_disk", map[string]string{"type": "pd-ssd"}, "pd-ssd"},
		{"google_storage_bucket", map[string]string{}, "standard"},
	}
	for _, tc := range tests {
		got := skuKey(tc.rtype, tc.attrs)
		if got != tc.want {
			t.Errorf("skuKey(%q, %v) = %q, want %q", tc.rtype, tc.attrs, got, tc.want)
		}
	}
}

func TestSkuKey_GenericFallback(t *testing.T) {
	got := skuKey("custom_resource", map[string]string{"sku": "my-sku"})
	if got != "my-sku" {
		t.Errorf("generic fallback: got %q, want %q", got, "my-sku")
	}
	got = skuKey("custom_resource", map[string]string{})
	if got != "" {
		t.Errorf("no sku attr: got %q, want empty", got)
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func writeTempPricingDB(t *testing.T, db PricingDB) string {
	t.Helper()
	data, err := json.Marshal(db)
	if err != nil {
		t.Fatalf("marshalling pricing db: %v", err)
	}
	path := filepath.Join(t.TempDir(), "pricing.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("writing temp pricing db: %v", err)
	}
	return path
}
