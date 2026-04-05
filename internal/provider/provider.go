// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

// Package provider implements the TerraKit Terraform/OpenTofu provider using
// the terraform-plugin-framework, ensuring compatibility with Terraform ≥ 1.5
// and OpenTofu ≥ 1.6.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/harpreet-x/terrakit/internal/pricing"
)

// Compile-time assertion: TerraKitProvider must implement provider.Provider.
var _ provider.Provider = &TerraKitProvider{}

// TerraKitProvider is the top-level provider struct.
type TerraKitProvider struct {
	version string
}

// TerraKitProviderModel mirrors the HCL schema for the provider block.
type TerraKitProviderModel struct {
	LocalPricingPath types.String `tfsdk:"local_pricing_path"`
}

// New returns a factory function that creates a new TerraKitProvider. The
// version string is injected at build time via -ldflags.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TerraKitProvider{version: version}
	}
}

// Metadata sets the provider type name and version reported to the framework.
func (p *TerraKitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "terrakit"
	resp.Version = p.version
}

// Schema defines the provider-level configuration attributes.
func (p *TerraKitProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
**TerraKit** — Privacy-First, Offline-Capable Terraform Cost Estimation.

Compute estimated infrastructure costs at plan-time without sending any data
to an external service. Point ` + "`local_pricing_path`" + ` at a local JSON pricing
database and every lookup happens entirely on your machine.
`,
		Attributes: map[string]schema.Attribute{
			"local_pricing_path": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Absolute or relative path to a local JSON pricing database. " +
					"When set, all cost lookups are resolved fully offline against this file. " +
					"When omitted, the provider falls back to its built-in illustrative pricing catalogue.",
			},
		},
	}
}

// Configure reads provider-level config, constructs the pricing engine once,
// and propagates it to data sources via ProviderData. Building the engine here
// means the JSON file is parsed exactly once per provider init and any file
// errors surface at configure-time rather than buried inside a Read call.
func (p *TerraKitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config TerraKitProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	engine, err := pricing.NewEngine(config.LocalPricingPath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Pricing Engine Initialisation Failed",
			"Could not load the pricing database: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = engine
	resp.ResourceData = engine
}

// DataSources returns the set of data sources exposed by this provider.
func (p *TerraKitProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCostDataSource,
	}
}

// Resources returns the set of managed resources. TerraKit is read-only by
// design; no resources are created in the remote cloud.
func (p *TerraKitProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
