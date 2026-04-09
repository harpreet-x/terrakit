// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

// terrakit is a privacy-first, offline-capable Terraform/OpenTofu/Terragrunt
// cost estimation CLI. It wraps your existing plan workflow and appends a
// fully-offline cost table — no data leaves your machine.
//
// Usage:
//
//	terrakit plan [terrakit-flags] [terraform/terragrunt flags]
//	terrakit cost [--plan-json=<path>] [--pricing-db=<path>]
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/harpreet-x/terrakit/internal/planparser"
	"github.com/harpreet-x/terrakit/internal/pricing"
)

const (
	// defaultPlanBin is the intermediate binary plan file written by `terraform plan -out`.
	defaultPlanBin = ".terrakit.tfplan"
	// defaultPlanJSON is the JSON plan file written by `terraform show -json` and
	// auto-detected by the terrakit_cost data source.
	defaultPlanJSON = ".terrakit.plan.json"

	divider = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd, args := os.Args[1], os.Args[2:]

	var err error
	switch cmd {
	case "plan":
		err = runPlan(args)
	case "cost":
		err = runCost(args)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "terrakit: unknown command %q\n\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "terrakit: %v\n", err)
		os.Exit(1)
	}
}

// ── plan subcommand ───────────────────────────────────────────────────────────

type planConfig struct {
	useTerragrunt bool
	useTerraform  bool
	pricingDB     string
	planOut       string
	workingDir    string
	passthrough   []string // forwarded verbatim to terraform/terragrunt
}

// parsePlanArgs separates TerraKit-specific flags from passthrough flags.
// All unrecognised args (terraform/terragrunt flags and positional args) are
// collected in passthrough and forwarded unchanged.
func parsePlanArgs(args []string) planConfig {
	c := planConfig{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--use-terragrunt":
			c.useTerragrunt = true
		case arg == "--use-terraform":
			c.useTerraform = true
		case arg == "--pricing-db" && i+1 < len(args):
			i++
			c.pricingDB = args[i]
		case strings.HasPrefix(arg, "--pricing-db="):
			c.pricingDB = strings.TrimPrefix(arg, "--pricing-db=")
		case arg == "--plan-out" && i+1 < len(args):
			i++
			c.planOut = args[i]
		case strings.HasPrefix(arg, "--plan-out="):
			c.planOut = strings.TrimPrefix(arg, "--plan-out=")
		case arg == "--working-dir" && i+1 < len(args):
			i++
			c.workingDir = args[i]
		case strings.HasPrefix(arg, "--working-dir="):
			c.workingDir = strings.TrimPrefix(arg, "--working-dir=")
		default:
			c.passthrough = append(c.passthrough, arg)
		}
	}
	return c
}

// runPlan implements `terrakit plan`:
//  1. Streams `terraform plan -out=<bin>` to the terminal.
//  2. Converts the binary plan to JSON with `terraform show -json`.
//  3. Prices the JSON plan offline and prints a cost table.
func runPlan(args []string) error {
	cfg := parsePlanArgs(args)

	if cfg.workingDir != "" {
		if err := os.Chdir(cfg.workingDir); err != nil {
			return fmt.Errorf("cannot change to working directory %q: %w", cfg.workingDir, err)
		}
	}

	planBin := cfg.planOut
	if planBin == "" {
		planBin = defaultPlanBin
	}

	tool, err := detectTool(cfg.useTerragrunt, cfg.useTerraform)
	if err != nil {
		return err
	}

	// ── Step 1: plan ─────────────────────────────────────────────────────────
	fmt.Println(divider)
	fmt.Printf(" [TerraKit] Step 1/2 — %s plan\n", tool)
	fmt.Println(divider)
	fmt.Println()

	planArgs := append([]string{"plan", "-out=" + planBin}, cfg.passthrough...)
	planCmd := exec.Command(tool, planArgs...)
	planCmd.Stdout = os.Stdout
	planCmd.Stderr = os.Stderr
	planCmd.Stdin = os.Stdin
	if err := planCmd.Run(); err != nil {
		return fmt.Errorf("%s plan: %w", tool, err)
	}

	// ── Step 2: cost estimate ─────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(divider)
	fmt.Println(" [TerraKit] Step 2/2 — cost estimate (offline)")
	fmt.Println(divider)
	fmt.Println()

	showOut, err := exec.Command(tool, "show", "-json", planBin).Output()
	if err != nil {
		return fmt.Errorf("%s show -json: %w", tool, err)
	}

	if err := os.WriteFile(defaultPlanJSON, showOut, 0o644); err != nil {
		return fmt.Errorf("writing plan JSON to %s: %w", defaultPlanJSON, err)
	}

	return printCostFromJSON(defaultPlanJSON, resolvePricingDB(cfg.pricingDB))
}

// ── cost subcommand ───────────────────────────────────────────────────────────

type costConfig struct {
	planJSON  string
	pricingDB string
}

func parseCostArgs(args []string) costConfig {
	c := costConfig{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--plan-json" && i+1 < len(args):
			i++
			c.planJSON = args[i]
		case strings.HasPrefix(arg, "--plan-json="):
			c.planJSON = strings.TrimPrefix(arg, "--plan-json=")
		case arg == "--pricing-db" && i+1 < len(args):
			i++
			c.pricingDB = args[i]
		case strings.HasPrefix(arg, "--pricing-db="):
			c.pricingDB = strings.TrimPrefix(arg, "--pricing-db=")
		}
	}
	return c
}

// runCost implements `terrakit cost`: prices an existing plan JSON without
// re-running a plan.
func runCost(args []string) error {
	cfg := parseCostArgs(args)

	planJSON := cfg.planJSON
	if planJSON == "" {
		planJSON = defaultPlanJSON
	}

	return printCostFromJSON(planJSON, resolvePricingDB(cfg.pricingDB))
}

// ── shared helpers ────────────────────────────────────────────────────────────

// resolvePricingDB returns the pricing DB path from the flag or the
// TERRAKIT_PRICING_DB environment variable (flag wins).
func resolvePricingDB(flag string) string {
	if flag != "" {
		return flag
	}
	return os.Getenv("TERRAKIT_PRICING_DB")
}

// detectTool returns "terraform" or "terragrunt" based on flags and the
// presence of a terragrunt.hcl file in the working directory.
func detectTool(forceGrunt, forceTF bool) (string, error) {
	if forceGrunt && forceTF {
		return "", fmt.Errorf("--use-terragrunt and --use-terraform are mutually exclusive")
	}
	if forceGrunt {
		return "terragrunt", nil
	}
	if forceTF {
		return "terraform", nil
	}
	if _, err := os.Stat("terragrunt.hcl"); err == nil {
		return "terragrunt", nil
	}
	return "terraform", nil
}

// printCostFromJSON parses the plan JSON at path, prices every resource with
// the given pricing DB (empty string = built-in catalogue), and prints the
// cost table to stdout.
func printCostFromJSON(path, pricingDB string) error {
	resources, err := planparser.ParseFile(path)
	if err != nil {
		return fmt.Errorf("parsing plan JSON %s: %w", path, err)
	}
	if len(resources) == 0 {
		fmt.Println("No resources to estimate (plan has no creates or updates).")
		return nil
	}

	engine, err := pricing.NewEngine(pricingDB)
	if err != nil {
		return fmt.Errorf("loading pricing engine: %w", err)
	}

	printCostTable(resources, engine)
	return nil
}

// ── cost table renderer ───────────────────────────────────────────────────────

func printCostTable(resources []planparser.Resource, engine *pricing.Engine) {
	const (
		colAddr    = 36
		colType    = 20
		colHourly  = 14
		colMonthly = 14
	)

	pad := func(s string, n int) string {
		if len(s) >= n {
			return s[:n-1] + " "
		}
		return s + strings.Repeat(" ", n-len(s))
	}
	money := func(f float64) string { return fmt.Sprintf("$%.4f", f) }

	top := "┌" + strings.Repeat("─", colAddr+2) + "┬" +
		strings.Repeat("─", colType+2) + "┬" +
		strings.Repeat("─", colHourly+2) + "┬" +
		strings.Repeat("─", colMonthly+2) + "┐"
	hr := "├" + strings.Repeat("─", colAddr+2) + "┼" +
		strings.Repeat("─", colType+2) + "┼" +
		strings.Repeat("─", colHourly+2) + "┼" +
		strings.Repeat("─", colMonthly+2) + "┤"
	bot := "└" + strings.Repeat("─", colAddr+2) + "┴" +
		strings.Repeat("─", colType+2) + "┴" +
		strings.Repeat("─", colHourly+2) + "┴" +
		strings.Repeat("─", colMonthly+2) + "┘"

	totalWidth := colAddr + colType + colHourly + colMonthly + 11
	title := "TerraKit Cost Estimate (offline)"
	titlePad := (totalWidth - len(title)) / 2

	fmt.Println(top)
	fmt.Printf("│%s%s%s│\n",
		strings.Repeat(" ", titlePad), title,
		strings.Repeat(" ", totalWidth-titlePad-len(title)))
	fmt.Println(hr)
	fmt.Printf("│ %s │ %s │ %s │ %s │\n",
		pad("Address", colAddr),
		pad("Type", colType),
		pad("Hourly (USD)", colHourly),
		pad("Monthly (USD)", colMonthly))
	fmt.Println(hr)

	var totalHourly, totalMonthly float64
	for _, r := range resources {
		est, err := engine.Estimate(r.Type, r.Attrs)
		if err != nil {
			fmt.Printf("│ %s │ %s │ %s │ %s │\n",
				pad(r.Address, colAddr),
				pad(r.Type, colType),
				pad("N/A", colHourly),
				pad("N/A", colMonthly))
			fmt.Fprintf(os.Stderr, "  [warn] %s: %v\n", r.Address, err)
			continue
		}
		totalHourly += est.HourlyUSD
		totalMonthly += est.MonthlyUSD
		fmt.Printf("│ %s │ %s │ %s │ %s │\n",
			pad(r.Address, colAddr),
			pad(r.Type, colType),
			pad(money(est.HourlyUSD), colHourly),
			pad(money(est.MonthlyUSD), colMonthly))
	}

	fmt.Println(hr)
	fmt.Printf("│ %s │ %s │ %s │ %s │\n",
		pad("TOTAL", colAddr),
		pad("", colType),
		pad(money(totalHourly), colHourly),
		pad(money(totalMonthly), colMonthly))
	fmt.Println(bot)
	fmt.Println("Prices: on-demand, us-east-1.  Hourly × 730 = Monthly.")
}

// ── usage ─────────────────────────────────────────────────────────────────────

func printUsage() {
	fmt.Print(`Usage: terrakit <command> [flags]

Commands:
  plan   Run terraform/terragrunt plan and print an offline cost estimate.
  cost   Print the cost table from an existing plan JSON (no re-plan).

terrakit plan [terrakit-flags] [terraform/terragrunt flags]
  --use-terragrunt          Force Terragrunt (default: auto-detect via terragrunt.hcl)
  --use-terraform           Force Terraform
  --pricing-db=<path>       Local JSON pricing database
  --plan-out=<path>         Override intermediate binary plan path
  --working-dir=<path>      Project directory (default: current directory)
  [all other flags are forwarded to terraform/terragrunt unchanged]

terrakit cost [flags]
  --plan-json=<path>        Plan JSON path (default: .terrakit.plan.json)
  --pricing-db=<path>       Local JSON pricing database

Environment:
  TERRAKIT_PRICING_DB       Pricing database path (overridden by --pricing-db)

Examples:
  terrakit plan
  terrakit plan -var="env=prod" -no-color
  terrakit plan --pricing-db=./pricing.json -var-file=prod.tfvars
  terrakit plan --working-dir=./envs/prod
  terrakit cost
  terrakit cost --plan-json=./ci-plan.json
`)
}
