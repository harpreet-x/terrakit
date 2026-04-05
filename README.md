# terrakit

**TerraKit** is a privacy-first, offline-capable cost estimation tool for
[Terraform](https://www.terraform.io/) (≥ 1.5), [OpenTofu](https://opentofu.org/) (≥ 1.6),
and [Terragrunt](https://terragrunt.gruntwork.io/).

> **Your infrastructure costs. Your machine. Your rules.**

---

## Why TerraKit?

### Privacy-First by Design

Most cloud-cost tools work by shipping your Terraform plan — containing resource
counts, instance types, region choices, and sometimes sensitive variable values —
to a third-party SaaS API. TerraKit inverts this model entirely:

| | Traditional SaaS cost tools | TerraKit |
|---|---|---|
| **Pricing lookups** | Remote API call | Local file / built-in DB |
| **Plan data leaves your machine** | Yes | **Never** |
| **Works in air-gapped environments** | No | **Yes** |
| **Custom / negotiated rates** | Vendor-specific UI | Drop a JSON file |
| **CI pipeline latency** | Network round-trip | Sub-millisecond |
| **Terragrunt support** | Varies | **Yes — auto-detected** |
| **Install required in project** | Often | **Never — install once globally** |

### Fully Offline

All pricing lookups run inside your own process — no sockets opened, no DNS
queries, no telemetry, no API keys required.

---

## Quick Start

### 1. Install the CLI once, globally

```bash
go install github.com/harpreet-x/terrakit/cmd/terrakit@latest
```

Or build from source:

```bash
git clone https://github.com/harpreet-x/terrakit
cd terrakit
make install     # builds terrakit binary and copies to $(GOPATH)/bin
```

> **You do not need to add any Makefile, script, or config file to your
> Terraform or Terragrunt project.** `terrakit` is a standalone binary
> installed once and used everywhere.

### 2. Use `terrakit plan` instead of `terraform plan` / `terragrunt plan`

Run it from your existing Terraform or Terragrunt project directory — no other
changes needed:

```bash
cd /your/terraform/project

# That's all — auto-detects terraform vs terragrunt
terrakit plan
```

**Every flag you already use passes through unchanged:**

```bash
# Terraform — all flags work
terrakit plan -var="env=prod"
terrakit plan -var-file=prod.tfvars
terrakit plan -target=aws_instance.web
terrakit plan -replace=aws_instance.web
terrakit plan -destroy
terrakit plan -refresh-only
terrakit plan -parallelism=30 -no-color
terrakit plan -lock=false -lock-timeout=30s
terrakit plan -out=myplan.tfplan         # honours your own -out

# Terragrunt — auto-detected when terragrunt.hcl exists
terrakit plan
terrakit plan -var="env=prod" --terragrunt-log-level=debug
terrakit plan --terragrunt-source-update
terrakit plan --working-dir=./envs/prod  # run from repo root

# Force one or the other
terrakit plan --use-terragrunt
terrakit plan --use-terraform
```

**Output — your normal plan output followed immediately by the cost table:**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 [TerraKit] Step 1/2 — terraform plan
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  # aws_instance.web will be created
  + resource "aws_instance" "web" { ... }
  # aws_instance.app will be created
  + resource "aws_instance" "app" { ... }

Plan: 2 to add, 0 to change, 0 to destroy.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 [TerraKit] Step 2/2 — cost estimate (offline)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

┌──────────────────────────────────────────────────────────────────────────────────────────┐
│                          TerraKit Cost Estimate (offline)                              │
├──────────────────────────────────────┬──────────────────────┬──────────────┬─────────────┤
│ Address                              │ Type                 │ Hourly (USD) │Monthly (USD)│
├──────────────────────────────────────┼──────────────────────┼──────────────┼─────────────┤
│ aws_instance.web                     │ aws_instance         │ $0.0104      │ $7.5920     │
│ aws_instance.app                     │ aws_instance         │ $0.0208      │ $15.1840    │
├──────────────────────────────────────┼──────────────────────┼──────────────┼─────────────┤
│ TOTAL                                │                      │ $0.0312      │ $22.7760    │
└──────────────────────────────────────┴──────────────────────┴──────────────┴─────────────┘
Prices: on-demand, us-east-1.  Hourly × 730 = Monthly.
```

---

## How It Works

Terraform's plugin protocol does not expose the full plan graph to providers
during plan execution. `terrakit plan` works around this transparently in two
steps — from your perspective it is a single command:

```
Step 1  terraform plan [your flags] -out=.terrakit.tfplan
        ↳ Full plan output streamed to your terminal. Binary plan saved locally.

Step 2  terraform show -json .terrakit.tfplan   (instant, no network)
        ↳ Binary plan converted to JSON.
        terrakit cost                            (offline, sub-millisecond)
        ↳ JSON parsed → prices looked up locally → cost table printed.
```

For Terragrunt, `terragrunt plan` and `terragrunt show` are used instead.
The plan JSON format is identical since Terragrunt delegates to Terraform.

---

## Terragrunt Support

TerraKit auto-detects Terragrunt projects — no configuration needed:

```
your-repo/
├── terragrunt.hcl          ← presence of this file triggers auto-detection
├── envs/
│   ├── prod/
│   │   └── terragrunt.hcl
│   └── staging/
│       └── terragrunt.hcl
└── modules/
    └── vpc/
        └── main.tf
```

```bash
# From within an env directory — auto-detected
cd envs/prod
terrakit plan

# From repo root targeting a specific environment
terrakit plan --working-dir=./envs/prod

# With Terragrunt-specific flags
terrakit plan --use-terragrunt --terragrunt-source-update -var="env=prod"
```

---

## Enforcing Budget Limits with `precondition`

Add the Terraform provider to your config to use `monthly_total` in
`lifecycle { precondition }` blocks. The plan fails with a full cost table in
the error message before any infrastructure is touched:

```hcl
# versions.tf
terraform {
  required_providers {
    terrakit = {
      source  = "registry.terraform.io/harpreet-x/terrakit"
      version = "~> 0.1"
    }
  }
}

provider "terrakit" {}

data "terrakit_cost" "estimate" {}   # auto-detects .terrakit.plan.json

resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  lifecycle {
    precondition {
      condition     = data.terrakit_cost.estimate.monthly_total < 100
      error_message = "Estimated monthly cost (${data.terrakit_cost.estimate.monthly_total} USD) exceeds the $100 budget cap.\n\n${data.terrakit_cost.estimate.summary}"
    }
  }
}
```

---

## `terrakit` CLI Reference

### `terrakit plan`

```
terrakit plan [terrakit flags] [terraform/terragrunt flags]

TerraKit flags:
  --use-terragrunt          Force terragrunt (default: auto-detect)
  --use-terraform           Force terraform
  --pricing-db=<path>       Local JSON pricing database
  --plan-out=<path>         Override intermediate plan binary path
  --working-dir=<path>      Project directory (default: current directory)
```

### `terrakit cost`

Re-print the cost table from an already-generated plan JSON without re-planning.
Useful in CI after `terrakit plan` has already run.

```bash
terrakit cost
terrakit cost --plan-json=./plan.json
terrakit cost --plan-json=./plan.json --pricing-db=./pricing.json
TERRAKIT_PRICING_DB=./pricing.json terrakit cost
```

---

## CI / CD Integration

Use `terrakit plan` in CI exactly as you would locally — no extra setup,
no secrets, no API tokens.

**GitHub Actions:**

```yaml
- name: Install terrakit
  run: go install github.com/harpreet-x/terrakit/cmd/terrakit@latest

- name: Plan + cost estimate
  run: terrakit plan -var-file=prod.tfvars -no-color
  env:
    TERRAKIT_PRICING_DB: ${{ github.workspace }}/pricing.json

- name: Upload plan JSON as artifact
  uses: actions/upload-artifact@v4
  with:
    name: plan-json
    path: .terrakit.plan.json
```

**GitLab CI:**

```yaml
plan:
  script:
    - go install github.com/harpreet-x/terrakit/cmd/terrakit@latest
    - terrakit plan -var-file=prod.tfvars -no-color -compact-warnings
  artifacts:
    paths:
      - .terrakit.plan.json
```

**Terragrunt mono-repo in CI:**

```yaml
- name: Plan all envs
  run: |
    for env in envs/*/; do
      echo "=== $env ==="
      terrakit plan --working-dir="$env" -no-color
    done
```

---

## `terrakit_cost` Terraform Data Source

### Input Arguments (choose one)

| Attribute | Type | Description |
|---|---|---|
| `plan_json_path` | `string` | Explicit path to a `terraform show -json` output file. |
| `resources` | `list(object)` | Manual resource list for ad-hoc estimates. |

When neither is set, auto-detects `.terrakit.plan.json` in the working
directory. When `plan_json_path` is set it takes priority over `resources`.

**`resources` object attributes:**

| Attribute | Required | Description |
|---|---|---|
| `name` | Yes | Logical label shown in the cost table. |
| `type` | Yes | Terraform resource type (e.g. `aws_instance`). |
| `attributes` | No | Key-value pairs used to select the pricing SKU. |

**Supported attribute keys by resource type:**

| Resource type | Key attribute | Example |
|---|---|---|
| `aws_instance` | `instance_type` | `"t3.micro"` |
| `aws_db_instance` | `instance_class` | `"db.t3.micro"` |
| Any | `sku` | any key in your pricing DB |

### Computed Attributes

| Attribute | Type | Description |
|---|---|---|
| `summary` | `string` | ASCII cost table. Use as an `output` or in a `precondition` error_message. |
| `line_items` | `list(object)` | Per-resource breakdown: `address`, `name`, `type`, `hourly_cost`, `monthly_cost`, `note`. |
| `monthly_total` | `number` | Aggregate monthly cost (hourly × 730). Safe in `precondition`. |
| `hourly_total` | `number` | Aggregate hourly cost. |
| `currency` | `string` | ISO 4217 code (currently always `"USD"`). |

---

## Local Pricing Database Format

Create a JSON file and pass it via `--pricing-db` (CLI) or `local_pricing_path`
(provider) to use your own rates — negotiated pricing, reserved instances, spot
estimates, or any custom values:

```json
{
  "aws_instance": {
    "t3.micro":  { "hourly_usd": 0.0104, "description": "t3.micro On-Demand Linux (us-east-1)" },
    "t3.small":  { "hourly_usd": 0.0208, "description": "t3.small On-Demand Linux (us-east-1)" },
    "t3.medium": { "hourly_usd": 0.0416, "description": "t3.medium On-Demand Linux (us-east-1)" },
    "t3.large":  { "hourly_usd": 0.0832, "description": "t3.large On-Demand Linux (us-east-1)" }
  },
  "aws_db_instance": {
    "db.t3.micro":  { "hourly_usd": 0.017, "description": "RDS MySQL db.t3.micro (us-east-1)" },
    "db.t3.small":  { "hourly_usd": 0.034, "description": "RDS MySQL db.t3.small (us-east-1)" }
  }
}
```

User-supplied entries are **merged on top of** the built-in catalogue — only
include the SKUs you want to override or add.

---

## Provider Configuration

| Attribute | Type | Required | Description |
|---|---|---|---|
| `local_pricing_path` | `string` | No | Path to a local JSON pricing database. When omitted, the built-in catalogue is used. |

---

## Known Limitation & Open Issues

Terraform's plugin protocol does not give providers access to the full plan
graph during plan execution. This is why `terrakit plan` uses a two-step
approach internally. Every cost estimation tool in this space works around this
the same way.

We have filed feature requests with both upstream projects:

- **Terraform:** [hashicorp/terraform#XXXXX](https://github.com/hashicorp/terraform/issues) —
  *Plugin protocol: expose planned resource graph to providers during plan phase*
- **OpenTofu:** [opentofu/opentofu#XXXXX](https://github.com/opentofu/opentofu/issues) —
  *RFC: provider plan context API — expose planned resource changes to plugins during plan*

A 👍 on those issues helps signal demand.

---

## Project Structure

```
terrakit/
├── main.go                          # Terraform provider entry point
├── go.mod
├── go.sum                           # Pinned dependency checksums (do not .gitignore)
├── GNUmakefile                      # build / install / test / vet targets
├── cmd/
│   └── terrakit/
│       └── main.go                  # terrakit CLI (plan + cost subcommands)
├── internal/
│   ├── provider/
│   │   ├── provider.go              # Provider schema & Configure
│   │   └── datasource_cost.go      # terrakit_cost data source
│   ├── pricing/
│   │   ├── engine.go               # Offline pricing engine (merge, lookup)
│   │   └── schema.go               # PricingDB types & built-in catalogue
│   └── planparser/
│       └── parser.go               # Terraform/Terragrunt plan JSON parser
└── examples/
    ├── basic/main.tf               # End-to-end HCL usage example
    └── pricing.json                # Sample local pricing database
```

---

## Building from Source

**Prerequisites:** Go 1.22+, GNU Make.

```bash
git clone https://github.com/harpreet-x/terrakit
cd terrakit
make build     # compiles both binaries into the repo root
make install   # installs terrakit into $(GOPATH)/bin and provider into plugin cache
```

**Available `make` targets:**

| Target | Description |
|---|---|
| `make build` | Compile provider binary + `terrakit` CLI. |
| `make install` | Build and install both into plugin cache / `$GOPATH/bin`. |
| `make test` | Run all unit tests. |
| `make vet` | Run `go vet` across all packages. |
| `make clean` | Remove compiled binaries. |

