// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

// Package planparser extracts priceable resources from a Terraform JSON plan
// produced by `terraform show -json <planfile>`. This is how TerraKit avoids
// requiring users to manually re-describe their infrastructure — it reads the
// plan Terraform already computed and prices what is actually going to change.
package planparser

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource is a single priceable resource extracted from the plan.
type Resource struct {
	// Address is the full Terraform resource address (e.g. "aws_instance.web").
	Address string
	// Name is the resource label (e.g. "web").
	Name string
	// Type is the resource type (e.g. "aws_instance").
	Type string
	// Attrs contains the string-valued planned attributes relevant for pricing
	// (e.g. {"instance_type": "t3.micro"}).
	Attrs map[string]string
}

// ── internal JSON shapes ────────────────────────────────────────────────────

type tfPlan struct {
	FormatVersion   string           `json:"format_version"`
	ResourceChanges []resourceChange `json:"resource_changes"`
}

type resourceChange struct {
	Address      string `json:"address"`
	ModuleAddress string `json:"module_address,omitempty"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Change       change `json:"change"`
}

type change struct {
	Actions []string               `json:"actions"`
	After   map[string]interface{} `json:"after"`
}

// ── Public API ───────────────────────────────────────────────────────────────

// ParseFile reads a Terraform JSON plan file and returns every resource that
// will be created or updated. Resources being destroyed or that have no-op
// changes are excluded — there is no cost impact from them.
//
// Generate the input file with:
//
//	terraform plan  -out=tfplan
//	terraform show  -json tfplan > plan.json
func ParseFile(path string) ([]Resource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading plan file %q: %w", path, err)
	}

	var plan tfPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("parsing plan JSON: %w", err)
	}

	var resources []Resource
	for _, rc := range plan.ResourceChanges {
		if !isChargeable(rc.Change.Actions) {
			continue
		}
		resources = append(resources, Resource{
			Address: rc.Address,
			Name:    rc.Name,
			Type:    rc.Type,
			Attrs:   flattenAttrs(rc.Change.After),
		})
	}

	return resources, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// isChargeable returns true when the action set means cloud resources will be
// provisioned or modified (create or update). Destroys and no-ops are free.
func isChargeable(actions []string) bool {
	for _, a := range actions {
		if a == "create" || a == "update" {
			return true
		}
	}
	return false
}

// flattenAttrs converts the plan's `after` map into map[string]string.
// Only string and numeric scalar values are included; nested objects/lists are
// skipped because pricing keys (instance_type, instance_class, etc.) are
// always top-level strings in the Terraform plan JSON.
func flattenAttrs(after map[string]interface{}) map[string]string {
	out := make(map[string]string, len(after))
	for k, v := range after {
		switch s := v.(type) {
		case string:
			out[k] = s
		case float64:
			// Numeric plan values (e.g. root_block_device.volume_size = 20)
			// are represented as float64 by encoding/json.
			out[k] = fmt.Sprintf("%g", s)
		case []interface{}:
			// List attributes (e.g. instance_types = ["m5.xlarge"]).
			// Extract the first string element so pricing lookups work.
			if len(s) > 0 {
				if str, ok := s[0].(string); ok {
					out[k] = str
				}
			}
		}
	}
	return out
}
