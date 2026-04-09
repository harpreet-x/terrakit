// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/harpreet-x/terrakit/internal/planparser"
	"github.com/harpreet-x/terrakit/internal/pricing"
)

// ── parsePlanArgs ─────────────────────────────────────────────────────────────

func TestParsePlanArgs_Defaults(t *testing.T) {
	cfg := parsePlanArgs([]string{})
	if cfg.useTerragrunt || cfg.useTerraform {
		t.Error("expected both tool flags false by default")
	}
	if cfg.pricingDB != "" || cfg.planOut != "" || cfg.workingDir != "" {
		t.Error("expected all string flags empty by default")
	}
	if len(cfg.passthrough) != 0 {
		t.Errorf("expected empty passthrough, got %v", cfg.passthrough)
	}
}

func TestParsePlanArgs_KnownFlags(t *testing.T) {
	cfg := parsePlanArgs([]string{
		"--use-terragrunt",
		"--pricing-db=/tmp/p.json",
		"--plan-out=/tmp/out.tfplan",
		"--working-dir=/srv/tf",
	})
	if !cfg.useTerragrunt {
		t.Error("expected useTerragrunt=true")
	}
	if cfg.pricingDB != "/tmp/p.json" {
		t.Errorf("pricingDB: got %q", cfg.pricingDB)
	}
	if cfg.planOut != "/tmp/out.tfplan" {
		t.Errorf("planOut: got %q", cfg.planOut)
	}
	if cfg.workingDir != "/srv/tf" {
		t.Errorf("workingDir: got %q", cfg.workingDir)
	}
}

func TestParsePlanArgs_SpaceSeparatedValues(t *testing.T) {
	cfg := parsePlanArgs([]string{"--pricing-db", "/tmp/p.json", "--plan-out", "/tmp/plan"})
	if cfg.pricingDB != "/tmp/p.json" {
		t.Errorf("pricingDB: got %q", cfg.pricingDB)
	}
	if cfg.planOut != "/tmp/plan" {
		t.Errorf("planOut: got %q", cfg.planOut)
	}
}

func TestParsePlanArgs_Passthrough(t *testing.T) {
	cfg := parsePlanArgs([]string{"-var=env=prod", "-no-color", "--pricing-db=x.json"})
	if len(cfg.passthrough) != 2 {
		t.Errorf("expected 2 passthrough args, got %d: %v", len(cfg.passthrough), cfg.passthrough)
	}
	if cfg.passthrough[0] != "-var=env=prod" {
		t.Errorf("passthrough[0]: got %q", cfg.passthrough[0])
	}
}

// ── parseCostArgs ─────────────────────────────────────────────────────────────

func TestParseCostArgs_Defaults(t *testing.T) {
	cfg := parseCostArgs([]string{})
	if cfg.planJSON != "" || cfg.pricingDB != "" {
		t.Error("expected empty defaults")
	}
}

func TestParseCostArgs_Flags(t *testing.T) {
	cfg := parseCostArgs([]string{"--plan-json=./plan.json", "--pricing-db=./db.json"})
	if cfg.planJSON != "./plan.json" {
		t.Errorf("planJSON: got %q", cfg.planJSON)
	}
	if cfg.pricingDB != "./db.json" {
		t.Errorf("pricingDB: got %q", cfg.pricingDB)
	}
}

func TestParseCostArgs_SpaceSeparated(t *testing.T) {
	cfg := parseCostArgs([]string{"--plan-json", "./plan.json"})
	if cfg.planJSON != "./plan.json" {
		t.Errorf("planJSON: got %q", cfg.planJSON)
	}
}

// ── detectTool ────────────────────────────────────────────────────────────────

func TestDetectTool_ForceTerraform(t *testing.T) {
	got, err := detectTool(false, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "terraform" {
		t.Errorf("got %q, want terraform", got)
	}
}

func TestDetectTool_ForceTerragrunt(t *testing.T) {
	got, err := detectTool(true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "terragrunt" {
		t.Errorf("got %q, want terragrunt", got)
	}
}

func TestDetectTool_BothForcedIsError(t *testing.T) {
	_, err := detectTool(true, true)
	if err == nil {
		t.Error("expected error when both force flags are set, got nil")
	}
}

func TestDetectTool_AutoDetectNoHCL(t *testing.T) {
	// Run from a temp dir with no terragrunt.hcl — expect terraform.
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	got, err := detectTool(false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "terraform" {
		t.Errorf("got %q, want terraform", got)
	}
}

func TestDetectTool_AutoDetectWithHCL(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "terragrunt.hcl"), []byte(""), 0o644)
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	got, err := detectTool(false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "terragrunt" {
		t.Errorf("got %q, want terragrunt", got)
	}
}

// ── resolvePricingDB ──────────────────────────────────────────────────────────

func TestResolvePricingDB_FlagWins(t *testing.T) {
	os.Setenv("TERRAKIT_PRICING_DB", "/env/path.json")
	defer os.Unsetenv("TERRAKIT_PRICING_DB")
	got := resolvePricingDB("/flag/path.json")
	if got != "/flag/path.json" {
		t.Errorf("got %q, want /flag/path.json", got)
	}
}

func TestResolvePricingDB_FallsBackToEnv(t *testing.T) {
	os.Setenv("TERRAKIT_PRICING_DB", "/env/path.json")
	defer os.Unsetenv("TERRAKIT_PRICING_DB")
	got := resolvePricingDB("")
	if got != "/env/path.json" {
		t.Errorf("got %q, want /env/path.json", got)
	}
}

func TestResolvePricingDB_EmptyWhenNeitherSet(t *testing.T) {
	os.Unsetenv("TERRAKIT_PRICING_DB")
	got := resolvePricingDB("")
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

// ── printCostTable (smoke test — no panic, produces expected output) ──────────

func TestPrintCostTable_Smoke(t *testing.T) {
	e, err := pricing.NewEngine("")
	if err != nil {
		t.Fatalf("engine: %v", err)
	}
	resources := []planparser.Resource{
		{Address: "aws_instance.web", Type: "aws_instance", Name: "web",
			Attrs: map[string]string{"instance_type": "t3.micro"}},
		{Address: "aws_nat_gateway.main", Type: "aws_nat_gateway", Name: "main",
			Attrs: map[string]string{}},
	}
	// Redirect stdout to capture output without printing during tests.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printCostTable(resources, e)

	w.Close()
	os.Stdout = old

	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		sb.Write(buf[:n])
		if err != nil {
			break
		}
	}
	out := sb.String()

	if !strings.Contains(out, "aws_instance.web") {
		t.Error("expected aws_instance.web in output")
	}
	if !strings.Contains(out, "TOTAL") {
		t.Error("expected TOTAL row in output")
	}
	if !strings.Contains(out, "$") {
		t.Error("expected dollar amounts in output")
	}
}

// ── printCostFromJSON (integration) ──────────────────────────────────────────

func TestPrintCostFromJSON_ValidPlan(t *testing.T) {
	planJSON := `{
		"format_version": "1.0",
		"resource_changes": [{
			"address": "aws_instance.test",
			"type": "aws_instance",
			"name": "test",
			"change": {
				"actions": ["create"],
				"after": { "instance_type": "t3.small" }
			}
		}]
	}`
	path := filepath.Join(t.TempDir(), "plan.json")
	os.WriteFile(path, []byte(planJSON), 0o644)

	// Redirect stdout to avoid polluting test output.
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	err := printCostFromJSON(path, "")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPrintCostFromJSON_EmptyPlan(t *testing.T) {
	planJSON := `{ "format_version": "1.0", "resource_changes": [] }`
	path := filepath.Join(t.TempDir(), "plan.json")
	os.WriteFile(path, []byte(planJSON), 0o644)

	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	err := printCostFromJSON(path, "")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPrintCostFromJSON_MissingFile(t *testing.T) {
	err := printCostFromJSON("/nonexistent/plan.json", "")
	if err == nil {
		t.Error("expected error for missing plan file, got nil")
	}
}
