terraform {
  required_providers {
    terrakit = {
      source  = "registry.terraform.io/harpreet-x/terrakit"
      version = "~> 0.1"
    }
  }
}

# ── Provider ──────────────────────────────────────────────────────────────────
provider "terrakit" {
  # Optional: point to a local JSON pricing database for custom/negotiated rates.
  # local_pricing_path = "${path.module}/../pricing.json"
}

# ── Auto-detect costs from plan (recommended) ─────────────────────────────────
#
# Step 1 — generate the plan JSON (run these in your shell before `terraform apply`):
#
#   terraform plan -out=tfplan
#   terraform show -json tfplan > plan.json
#
# Step 2 — TerraKit reads plan.json, prices every create/update automatically.
#
data "terrakit_cost" "auto" {
  plan_json_path = "${path.module}/plan.json"
}

# Print the full cost table to the terminal.
output "cost_summary" {
  value = data.terrakit_cost.auto.summary
}

output "monthly_total_usd" {
  value = data.terrakit_cost.auto.monthly_total
}

# Per-resource breakdown — useful for CI cost reports.
output "line_items" {
  value = data.terrakit_cost.auto.line_items
}

# ── Budget guard via precondition ─────────────────────────────────────────────
#
# Any resource can block apply when the estimate exceeds a budget cap.
# The plan fails BEFORE any infrastructure is touched.
#
# resource "aws_instance" "web" {
#   ami           = "ami-0c55b159cbfafe1f0"
#   instance_type = "t3.micro"
#
#   lifecycle {
#     precondition {
#       condition     = data.terrakit_cost.auto.monthly_total < 100
#       error_message = "Estimated monthly cost (${data.terrakit_cost.auto.monthly_total} USD) exceeds the $100 budget cap.\n\n${data.terrakit_cost.auto.summary}"
#     }
#   }
# }

# ── Manual list (ad-hoc estimates without a plan file) ────────────────────────
data "terrakit_cost" "manual" {
  resources = [
    {
      name       = "web_server"
      type       = "aws_instance"
      attributes = { instance_type = "t3.micro" }
    },
    {
      name       = "app_server"
      type       = "aws_instance"
      attributes = { instance_type = "t3.small" }
    },
  ]
}

output "manual_cost_summary" {
  value = data.terrakit_cost.manual.summary
}
