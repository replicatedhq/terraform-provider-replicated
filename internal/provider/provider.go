package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/replicatedhq/replicated/pkg/kotsclient"
	"github.com/replicatedhq/replicated/pkg/platformclient"
)

// Ensure ReplicatedProvider satisfies various provider interfaces.
var _ provider.Provider = &ReplicatedProvider{}

// ReplicatedProvider defines the provider implementation.
type ReplicatedProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ReplicatedProviderModel describes the provider data model.
type ReplicatedProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiToken types.String `tfsdk:"api_token"`
}

func (p *ReplicatedProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "replicated"
	resp.Version = p.version
}

func (p *ReplicatedProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Vendor API endpoint",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Vendor API token",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ReplicatedProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiOrigin := os.Getenv("REPLICATED_API_ORIGIN")
	apiToken := os.Getenv("REPLICATED_API_TOKEN")
	var data ReplicatedProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Endpoint.ValueString() != "" {
		apiOrigin = data.Endpoint.ValueString()
	}

	if data.ApiToken.ValueString() != "" {
		apiToken = data.ApiToken.ValueString()
	}

	if apiOrigin == "" {
		apiOrigin = "https://api.replicated.com/vendor"
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				"the REPLICATED_API_TOKEN environment variable or provider "+
				"configuration block api_token attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpClient := platformclient.NewHTTPClient(apiOrigin, apiToken)
	kotsAPI := &kotsclient.VendorV3Client{HTTPClient: *httpClient}

	// VendorApi client configuration.
	client := kotsAPI
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ReplicatedProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
	}
}

func (p *ReplicatedProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ReplicatedProvider{
			version: version,
		}
	}
}
