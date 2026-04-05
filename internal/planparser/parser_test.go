// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

package planparser

import (
	"os"
	"path/filepath"
	"testing"
)

// ── isChargeable ─────────────────────────────────────────────────────────────

func TestIsChargeable(t *testing.T) {
	tests := []struct {
		actions []string
		want    bool
	}{
		{[]string{"create"}, true},
		{[]string{"update"}, true},
		{[]string{"create", "update"}, true},
		{[]string{"delete"}, false},
		{[]string{"no-op"}, false},
		{[]string{}, false},
		{[]string{"delete", "create"}, true}, // replace — still costs money
	}
	for _, tc := range tests {
		got := isChargeable(tc.actions)
		if got != tc.want {
			t.Errorf("isChargeable(%v) = %v, want %v", tc.actions, got, tc.want)
		}
	}
}

// ── flattenAttrs ─────────────────────────────────────────────────────────────

func TestFlattenAttrs(t *testing.T) {
	in := map[string]interface{}{
		"instance_type":  "t3.micro",
		"volume_size":    float64(20),
		"instance_types": []interface{}{"m5.xlarge", "m5.2xlarge"},
		"tags":           map[string]interface{}{"env": "prod"}, // nested — skipped
		"count":          nil,                                   // nil — skipped
	}
	got := flattenAttrs(in)

	assertEqual(t, "instance_type", got["instance_type"], "t3.micro")
	assertEqual(t, "volume_size", got["volume_size"], "20")
	assertEqual(t, "instance_types", got["instance_types"], "m5.xlarge") // first element
	if _, ok := got["tags"]; ok {
		t.Error("nested map should be skipped")
	}
}

func TestFlattenAttrsEmpty(t *testing.T) {
	got := flattenAttrs(nil)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

// ── ParseFile ────────────────────────────────────────────────────────────────

func TestParseFile_CreateAndUpdate(t *testing.T) {
	raw := `{
		"format_version": "1.0",
		"resource_changes": [
			{
				"address": "aws_instance.web",
				"type": "aws_instance",
				"name": "web",
				"change": { "actions": ["create"], "after": { "instance_type": "t3.micro" } }
			},
			{
				"address": "aws_instance.db",
				"type": "aws_instance",
				"name": "db",
				"change": { "actions": ["update"], "after": { "instance_type": "t3.large" } }
			}
		]
	}`
	path := writeTempJSON(t, raw)
	resources, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}
	assertEqual(t, "address[0]", resources[0].Address, "aws_instance.web")
	assertEqual(t, "attrs[0]", resources[0].Attrs["instance_type"], "t3.micro")
	assertEqual(t, "address[1]", resources[1].Address, "aws_instance.db")
	assertEqual(t, "attrs[1]", resources[1].Attrs["instance_type"], "t3.large")
}

func TestParseFile_SkipsDestroyAndNoOp(t *testing.T) {
	raw := `{
		"format_version": "1.0",
		"resource_changes": [
			{
				"address": "aws_instance.old",
				"type": "aws_instance",
				"name": "old",
				"change": { "actions": ["delete"], "after": {} }
			},
			{
				"address": "aws_instance.same",
				"type": "aws_instance",
				"name": "same",
				"change": { "actions": ["no-op"], "after": { "instance_type": "t3.micro" } }
			}
		]
	}`
	path := writeTempJSON(t, raw)
	resources, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 0 {
		t.Errorf("expected 0 resources, got %d: %+v", len(resources), resources)
	}
}

func TestParseFile_EmptyPlan(t *testing.T) {
	raw := `{ "format_version": "1.0", "resource_changes": [] }`
	path := writeTempJSON(t, raw)
	resources, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(resources))
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/plan.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseFile_InvalidJSON(t *testing.T) {
	path := writeTempJSON(t, `{ not valid json }`)
	_, err := ParseFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func writeTempJSON(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "plan.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func assertEqual(t *testing.T, label, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", label, got, want)
	}
}
