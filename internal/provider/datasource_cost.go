// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

package provider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/harpreet-x/terrakit/internal/planparser"
	"github.com/harpreet-x/terrakit/internal/pricing"
)

// autoDetectPlanFile is the filename `make plan` writes. When it exists in the
// working directory and plan_json_path is not explicitly configured, TerraKit
// uses it automatically so users get cost output with zero extra config.
const autoDetectPlanFile = ".terrakit.plan.json"

// Compile-time assertions.
var _ datasource.DataSource = &CostDataSource{}
var _ datasource.DataSourceWithConfigure = &CostDataSource{}

// ─── Model types ─────────────────────────────────────────────────────────────

// ResourceItemModel is used when the caller provides a manual resource list.
type ResourceItemModel struct {
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Attributes types.Map    `tfsdk:"attributes"`
}

// LineItemModel is one row in the computed cost breakdown.
type LineItemModel struct {
	Name        types.String  `tfsdk:"name"`
	Address     types.String  `tfsdk:"address"`
	Type        types.String  `tfsdk:"type"`
	HourlyCost  types.Float64 `tfsdk:"hourly_cost"`
	MonthlyCost types.Float64 `tfsdk:"monthly_cost"`
	Note        types.String  `tfsdk:"note"`
}

// CostDataSourceModel is the full state model for the terrakit_cost data source.
type CostDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	// ── Input (choose one) ───────────────────────────────────────────────────

	// PlanJSONPath: path to `terraform show -json <planfile>` output.
	// When set, resources are auto-detected from the plan — no manual list needed.
	PlanJSONPath types.String `tfsdk:"plan_json_path"`

	// Resources: explicit list for ad-hoc estimates without a plan file.
	Resources []ResourceItemModel `tfsdk:"resources"`

	// ── Computed outputs ─────────────────────────────────────────────────────

	LineItems    []LineItemModel `tfsdk:"line_items"`
	MonthlyTotal types.Float64   `tfsdk:"monthly_total"`
	HourlyTotal  types.Float64   `tfsdk:"hourly_total"`
	Currency     types.String    `tfsdk:"currency"`
	// Summary is a formatted ASCII cost table — pipe it to an output or use it
	// in a precondition error_message for a human-readable budget failure.
	Summary types.String `tfsdk:"summary"`
}

// ─── Constructor ─────────────────────────────────────────────────────────────

type CostDataSource struct {
	engine *pricing.Engine
}

func NewCostDataSource() datasource.DataSource {
	return &CostDataSource{}
}

// ─── Metadata ────────────────────────────────────────────────────────────────

func (d *CostDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cost"
}

// ─── Schema ──────────────────────────────────────────────────────────────────

func (d *CostDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	lineItemAttrs := map[string]schema.Attribute{
		"name":         schema.StringAttribute{Computed: true, MarkdownDescription: "Resource label."},
		"address":      schema.StringAttribute{Computed: true, MarkdownDescription: "Full Terraform resource address (e.g. `aws_instance.web`)."},
		"type":         schema.StringAttribute{Computed: true, MarkdownDescription: "Terraform resource type."},
		"hourly_cost":  schema.Float64Attribute{Computed: true, MarkdownDescription: "Estimated hourly cost in `currency`."},
		"monthly_cost": schema.Float64Attribute{Computed: true, MarkdownDescription: "Estimated monthly cost in `currency` (hourly × 730)."},
		"note":         schema.StringAttribute{Computed: true, MarkdownDescription: "Empty when priced successfully; contains a warning when the SKU was not found."},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: `
Estimates the cost of infrastructure resources — either auto-detected from a
Terraform plan file or specified manually.

**Auto-detect from plan (recommended):**

` + "```bash" + `
terraform plan  -out=tfplan
terraform show  -json tfplan > plan.json
` + "```" + `

` + "```hcl" + `
data "terrakit_cost" "estimate" {
  plan_json_path = "${path.module}/plan.json"
}

output "cost_summary" {
  value = data.terrakit_cost.estimate.summary
}
` + "```" + `

All lookups are fully offline — no data leaves your machine.`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},

			// ── Inputs ───────────────────────────────────────────────────────

			"plan_json_path": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Path to a Terraform JSON plan file produced by " +
					"`terraform show -json <planfile>`. When set, resources are auto-detected " +
					"from the plan — no manual `resources` list is needed.",
			},

			"resources": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Manual resource list for ad-hoc estimates. Ignored when `plan_json_path` is set.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Logical label (used in the summary table and diagnostics).",
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Terraform resource type (e.g. `aws_instance`).",
						},
						"attributes": schema.MapAttribute{
							Optional:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "Attributes used to select the pricing SKU (e.g. `{ instance_type = \"t3.micro\" }`).",
						},
					},
				},
			},

			// ── Computed outputs ─────────────────────────────────────────────

			"line_items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Per-resource cost breakdown. One entry per resource from the plan or manual list.",
				NestedObject:        schema.NestedAttributeObject{Attributes: lineItemAttrs},
			},

			"monthly_total": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "Aggregate estimated monthly cost. Safe to use in `lifecycle { precondition }` budget guards.",
			},

			"hourly_total": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "Aggregate estimated hourly cost.",
			},

			"currency": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ISO 4217 currency code (currently always `\"USD\"`).",
			},

			"summary": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable ASCII cost table. Use as an `output` value or inside a `precondition` error_message.",
			},
		},
	}
}

// ─── Configure ───────────────────────────────────────────────────────────────

func (d *CostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	engine, ok := req.ProviderData.(*pricing.Engine)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *pricing.Engine, got %T.", req.ProviderData),
		)
		return
	}
	d.engine = engine
}

// ─── Read ─────────────────────────────────────────────────────────────────────

func (d *CostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CostDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve the resource list from the plan file or the manual list.
	resources, diags := d.resolveResources(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(resources) == 0 {
		resp.Diagnostics.AddError(
			"No Resources to Estimate",
			"Provide either plan_json_path (pointing to `terraform show -json` output) "+
				"or a non-empty resources list.",
		)
		return
	}

	// Price every resource.
	lineItems, totalHourly, totalMonthly := d.priceResources(resources, resp)

	// Build a stable content-based ID.
	var idBuf strings.Builder
	for _, r := range resources {
		idBuf.WriteString(r.Address)
		idBuf.WriteByte('|')
	}
	sum := sha256.Sum256([]byte(idBuf.String()))

	state.ID = types.StringValue(fmt.Sprintf("terrakit-%x", sum[:8]))
	state.LineItems = lineItems
	state.HourlyTotal = types.Float64Value(totalHourly)
	state.MonthlyTotal = types.Float64Value(totalMonthly)
	state.Currency = types.StringValue("USD")
	state.Summary = types.StringValue(buildSummary(lineItems, totalHourly, totalMonthly))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// ─── resolveResources ────────────────────────────────────────────────────────

// internalResource is a normalised resource record used internally regardless
// of whether it came from a plan file or a manual list.
type internalResource struct {
	Address string
	Name    string
	Type    string
	Attrs   map[string]string
}

func (d *CostDataSource) resolveResources(ctx context.Context, state *CostDataSourceModel) ([]internalResource, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Resolve the plan file path: explicit config → auto-detect → manual list.
	planPath := state.PlanJSONPath.ValueString()
	if planPath == "" {
		if _, err := os.Stat(autoDetectPlanFile); err == nil {
			planPath = autoDetectPlanFile
		}
	}

	if planPath != "" {
		parsed, err := planparser.ParseFile(planPath)
		if err != nil {
			diags.AddError("Plan File Parse Error", err.Error())
			return nil, diags
		}
		out := make([]internalResource, len(parsed))
		for i, r := range parsed {
			out[i] = internalResource{
				Address: r.Address,
				Name:    r.Name,
				Type:    r.Type,
				Attrs:   r.Attrs,
			}
		}
		return out, diags
	}

	// Fall back to manual resources list.
	out := make([]internalResource, 0, len(state.Resources))
	for _, r := range state.Resources {
		attrs, attrDiags := flattenAttributes(ctx, r.Attributes)
		diags.Append(attrDiags...)
		if diags.HasError() {
			return nil, diags
		}
		name := r.Name.ValueString()
		rtype := r.Type.ValueString()
		out = append(out, internalResource{
			Address: rtype + "." + name,
			Name:    name,
			Type:    rtype,
			Attrs:   attrs,
		})
	}
	return out, diags
}

// ─── priceResources ──────────────────────────────────────────────────────────

func (d *CostDataSource) priceResources(resources []internalResource, resp *datasource.ReadResponse) ([]LineItemModel, float64, float64) {
	var totalHourly, totalMonthly float64
	lineItems := make([]LineItemModel, 0, len(resources))

	for _, r := range resources {
		estimate, err := d.engine.Estimate(r.Type, r.Attrs)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Pricing Lookup Skipped",
				fmt.Sprintf("%s (%s): %s", r.Address, r.Type, err.Error()),
			)
			lineItems = append(lineItems, LineItemModel{
				Name:        types.StringValue(r.Name),
				Address:     types.StringValue(r.Address),
				Type:        types.StringValue(r.Type),
				HourlyCost:  types.Float64Value(0),
				MonthlyCost: types.Float64Value(0),
				Note:        types.StringValue("SKU not found: " + err.Error()),
			})
			continue
		}

		totalHourly += estimate.HourlyUSD
		totalMonthly += estimate.MonthlyUSD
		lineItems = append(lineItems, LineItemModel{
			Name:        types.StringValue(r.Name),
			Address:     types.StringValue(r.Address),
			Type:        types.StringValue(r.Type),
			HourlyCost:  types.Float64Value(estimate.HourlyUSD),
			MonthlyCost: types.Float64Value(estimate.MonthlyUSD),
			Note:        types.StringValue(""),
		})
	}

	return lineItems, totalHourly, totalMonthly
}

// ─── buildSummary ────────────────────────────────────────────────────────────

// buildSummary renders a fixed-width ASCII cost table suitable for terminal
// output, Terraform plan output, and precondition error messages.
//
// Example:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                   TerraKit Cost Estimate                      │
//	├───────────────────────┬──────────────┬────────────┬────────────┤
//	│ Address               │ Type         │ Hourly     │ Monthly    │
//	├───────────────────────┼──────────────┼────────────┼────────────┤
//	│ aws_instance.web      │ aws_instance │ $0.0104    │ $7.59      │
//	│ aws_instance.db       │ aws_instance │ $0.0208    │ $15.18     │
//	├───────────────────────┼──────────────┼────────────┼────────────┤
//	│ TOTAL                 │              │ $0.0312    │ $22.78     │
//	└───────────────────────┴──────────────┴────────────┴────────────┘
func buildSummary(items []LineItemModel, totalHourly, totalMonthly float64) string {
	const (
		colAddr    = 28
		colType    = 16
		colHourly  = 13 // "Hourly" = 6 chars; "$0.0000" = 7 chars — 13 fits both cleanly
		colMonthly = 13 // "Monthly" = 7 chars; "$0.0000" = 7 chars — 13 fits both cleanly
	)

	pad := func(s string, n int) string {
		if len(s) >= n {
			return s[:n-1] + " "
		}
		return s + strings.Repeat(" ", n-len(s))
	}

	money := func(f float64) string { return fmt.Sprintf("$%.4f", f) }

	hr := "├" + strings.Repeat("─", colAddr+2) + "┼" +
		strings.Repeat("─", colType+2) + "┼" +
		strings.Repeat("─", colHourly+2) + "┼" +
		strings.Repeat("─", colMonthly+2) + "┤"

	totalWidth := colAddr + colType + colHourly + colMonthly + 11 // borders + spaces
	title := "TerraKit Cost Estimate"
	titlePad := (totalWidth - len(title)) / 2

	var sb strings.Builder

	// Top border + title
	sb.WriteString("┌" + strings.Repeat("─", totalWidth) + "┐\n")
	sb.WriteString("│" + strings.Repeat(" ", titlePad) + title +
		strings.Repeat(" ", totalWidth-titlePad-len(title)) + "│\n")

	// Header row
	sb.WriteString(hr + "\n")
	sb.WriteString("│ " + pad("Address", colAddr) +
		" │ " + pad("Type", colType) +
		" │ " + pad("Hourly", colHourly) +
		" │ " + pad("Monthly", colMonthly) + " │\n")
	sb.WriteString(hr + "\n")

	// Data rows
	for _, item := range items {
		note := item.Note.ValueString()
		hourlyCost := money(item.HourlyCost.ValueFloat64())
		monthlyCost := money(item.MonthlyCost.ValueFloat64())
		if note != "" {
			hourlyCost = "N/A"
			monthlyCost = "N/A"
		}
		sb.WriteString("│ " + pad(item.Address.ValueString(), colAddr) +
			" │ " + pad(item.Type.ValueString(), colType) +
			" │ " + pad(hourlyCost, colHourly) +
			" │ " + pad(monthlyCost, colMonthly) + " │\n")
	}

	// Total row
	sb.WriteString(hr + "\n")
	sb.WriteString("│ " + pad("TOTAL", colAddr) +
		" │ " + pad("", colType) +
		" │ " + pad(money(totalHourly), colHourly) +
		" │ " + pad(money(totalMonthly), colMonthly) + " │\n")

	// Bottom border
	sb.WriteString("└" + strings.Repeat("─", colAddr+2) + "┴" +
		strings.Repeat("─", colType+2) + "┴" +
		strings.Repeat("─", colHourly+2) + "┴" +
		strings.Repeat("─", colMonthly+2) + "┘\n")

	sb.WriteString("Currency: USD  |  Hourly × 730 = Monthly")

	return sb.String()
}

// ─── flattenAttributes ───────────────────────────────────────────────────────

func flattenAttributes(ctx context.Context, m types.Map) (map[string]string, diag.Diagnostics) {
	out := make(map[string]string)
	if m.IsNull() || m.IsUnknown() {
		return out, nil
	}
	elements := make(map[string]attr.Value)
	diags := m.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return out, diags
	}
	for k, v := range elements {
		if sv, ok := v.(types.String); ok {
			out[k] = sv.ValueString()
		}
	}
	return out, nil
}
