package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/replicatedhq/replicated/pkg/kotsclient"
)

var _ resource.Resource = &CustomerResource{}
var _ resource.ResourceWithImportState = &CustomerResource{}

func NewCustomerResource() resource.Resource {
	return &CustomerResource{}
}

type CustomerResource struct {
	client     *CustomersAPIClient
	kotsClient *kotsclient.VendorV3Client
}

type CustomerResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	AppId                            types.String `tfsdk:"app_id"`
	ChannelId                        types.String `tfsdk:"channel_id"`
	Email                            types.String `tfsdk:"email"`
	EntitlementValues                types.Map    `tfsdk:"entitlement_values"`
	ExpiresAt                        types.String `tfsdk:"expires_at"`
	IsAirgapEnabled                  types.Bool   `tfsdk:"is_airgap_enabled"`
	IsEmbeddedClusterDownloadEnabled types.Bool   `tfsdk:"is_embedded_cluster_download_enabled"`
	IsGeoaxisSupported               types.Bool   `tfsdk:"is_geoaxis_supported"`
	IsGitopsSupported                types.Bool   `tfsdk:"is_gitops_supported"`
	IsHelmvmDownloadEnabled          types.Bool   `tfsdk:"is_helmvm_download_enabled"`
	IsIdentityServiceSupported       types.Bool   `tfsdk:"is_identity_service_supported"`
	IsInstallerSupportEnabled        types.Bool   `tfsdk:"is_installer_support_enabled"`
	IsKotsInstallEnabled             types.Bool   `tfsdk:"is_kots_install_enabled"`
	IsSnapshotSupported              types.Bool   `tfsdk:"is_snapshot_supported"`
	IsSupportBundleUploadEnabled     types.Bool   `tfsdk:"is_support_bundle_upload_enabled"`
	Name                             types.String `tfsdk:"name"`
	Type                             types.String `tfsdk:"type"`
}

func (r *CustomerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customer"
}

func (r *CustomerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Customer resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the customer",
				Computed:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "App to which the channel is associated",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email of the customer",
				Optional:            true,
			},
			"entitlement_values": schema.MapAttribute{
				MarkdownDescription: "Entitlement values of the customer",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"expires_at": schema.StringAttribute{
				MarkdownDescription: "Expiration date of the customer license",
				Optional:            true,
			},
			"is_airgap_enabled": schema.BoolAttribute{
				MarkdownDescription: "Is airgap enabled for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_embedded_cluster_download_enabled": schema.BoolAttribute{
				MarkdownDescription: "Is embedded cluster download enabled for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_geoaxis_supported": schema.BoolAttribute{
				MarkdownDescription: "Is geoaxis supported for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_gitops_supported": schema.BoolAttribute{
				MarkdownDescription: "Is gitops supported for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_helmvm_download_enabled": schema.BoolAttribute{
				MarkdownDescription: "Is helmvm download enabled for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_identity_service_supported": schema.BoolAttribute{
				MarkdownDescription: "Is identity service supported for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_installer_support_enabled": schema.BoolAttribute{
				MarkdownDescription: "Is installer support enabled for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_kots_install_enabled": schema.BoolAttribute{
				MarkdownDescription: "Is kots install enabled for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_snapshot_supported": schema.BoolAttribute{
				MarkdownDescription: "Is snapshot supported for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_support_bundle_upload_enabled": schema.BoolAttribute{
				MarkdownDescription: "Is support bundle upload enabled for the customer license",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"channel_id": schema.StringAttribute{
				MarkdownDescription: "Channel to which the customer license is associated",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the customer",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the customer",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("trial"),
			},
		},
	}
}

func (r *CustomerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clients, ok := req.ProviderData.(*ReplicatedProviderClients)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &clients.customersAPIClient
	r.kotsClient = &clients.kotsVendorV3Client
}

func (r *CustomerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CustomerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var entitlementValues []EntitlementValue

	for name, value := range data.EntitlementValues.Elements() {
		entitlementValues = append(entitlementValues, EntitlementValue{Name: name, Value: strings.Trim(value.String(), "\"")})
	}

	opts := CreateCustomerOpts{
		AppID:                            data.AppId.ValueString(),
		ChannelId:                        data.ChannelId.ValueString(),
		Email:                            data.Email.ValueString(),
		EntitlementValues:                entitlementValues,
		ExpiresAt:                        data.ExpiresAt.ValueString(),
		IsAirgapEnabled:                  data.IsAirgapEnabled.ValueBool(),
		IsEmbeddedClusterDownloadEnabled: data.IsEmbeddedClusterDownloadEnabled.ValueBool(),
		IsGeoaxisSupported:               data.IsGeoaxisSupported.ValueBool(),
		IsGitopsSupported:                data.IsGitopsSupported.ValueBool(),
		IsHelmVMDownloadEnabled:          data.IsHelmvmDownloadEnabled.ValueBool(),
		IsIdentityServiceSupported:       data.IsIdentityServiceSupported.ValueBool(),
		IsInstallerSupportEnabled:        data.IsInstallerSupportEnabled.ValueBool(),
		IsKotsInstallEnabled:             data.IsKotsInstallEnabled.ValueBool(),
		IsSnapshotSupported:              data.IsSnapshotSupported.ValueBool(),
		IsSupportBundleUploadEnabled:     data.IsSupportBundleUploadEnabled.ValueBool(),
		Name:                             data.Name.ValueString(),
		Type:                             data.Type.ValueString(),
	}

	customer, err := r.client.CreateCustomer(opts)
	if err != nil {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to create customer, got error: %s", err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("app/%s/customer/%s", customer.Customer.Channels[0].AppID, customer.Customer.ID))

	tflog.Trace(ctx, "created a customer")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CustomerResourceModel

	var resourceId string

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &resourceId)...)

	if resp.Diagnostics.HasError() {
		return
	}

	appId := strings.Split(resourceId, "/")[1]
	id := strings.Split(resourceId, "/")[3]

	customer, err := r.client.GetCustomer(appId, id)
	if err != nil {
		if err.Error() == "Customer not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to get customer, got error: %s", err))
		return
	}

	entitlements := make(map[string]string)

	for _, entitlement := range customer.Customer.Entitlements {
		entitlements[entitlement.Name] = entitlement.Value
	}

	mapValue, diags := types.MapValueFrom(ctx, types.StringType, entitlements)

	if diags.HasError() {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to get customer, got error: %s", diags))
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("app/%s/customer/%s", customer.Customer.Channels[0].AppID, customer.Customer.ID))
	data.AppId = types.StringValue(customer.Customer.Channels[0].AppID)
	data.ChannelId = types.StringValue(customer.Customer.Channels[0].ID)
	data.Email = types.StringValue(customer.Customer.Email)
	data.EntitlementValues = mapValue
	data.ExpiresAt = types.StringValue(customer.Customer.ExpiresAt)
	data.IsAirgapEnabled = types.BoolValue(customer.Customer.Airgap)
	data.IsEmbeddedClusterDownloadEnabled = types.BoolValue(customer.Customer.IsEmbeddedClusterDownloadEnabled)
	data.IsGeoaxisSupported = types.BoolValue(customer.Customer.IsGeoaxisSupported)
	data.IsHelmvmDownloadEnabled = types.BoolValue(customer.Customer.IsHelmVMDownloadEnabled)
	data.IsIdentityServiceSupported = types.BoolValue(customer.Customer.IsIdentityServiceSupported)
	data.IsInstallerSupportEnabled = types.BoolValue(customer.Customer.IsInstallerSupportEnabled)
	data.IsKotsInstallEnabled = types.BoolValue(customer.Customer.IsKotsInstallEnabled)
	data.IsSnapshotSupported = types.BoolValue(customer.Customer.IsSnapshotSupported)
	data.IsSupportBundleUploadEnabled = types.BoolValue(customer.Customer.IsSupportBundleUploadEnabled)
	data.IsGitopsSupported = types.BoolValue(customer.Customer.IsGitopsSupported)
	data.Name = types.StringValue(customer.Customer.Name)
	data.Type = types.StringValue(customer.Customer.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var oldData CustomerResourceModel
	var updatedData CustomerResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &oldData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &updatedData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var entitlementValues []EntitlementValue

	for name, value := range updatedData.EntitlementValues.Elements() {
		entitlementValues = append(entitlementValues, EntitlementValue{Name: name, Value: strings.Trim(value.String(), "\"")})
	}

	var opts UpdateCustomerOpts

	opts.AppID = updatedData.AppId.ValueString()
	opts.ChannelId = updatedData.ChannelId.ValueString()
	opts.Email = updatedData.Email.ValueString()
	opts.EntitlementValues = entitlementValues
	opts.ExpiresAt = updatedData.ExpiresAt.ValueString()
	opts.IsAirgapEnabled = updatedData.IsAirgapEnabled.ValueBool()
	opts.IsEmbeddedClusterDownloadEnabled = updatedData.IsEmbeddedClusterDownloadEnabled.ValueBool()
	opts.IsGeoaxisSupported = updatedData.IsGeoaxisSupported.ValueBool()
	opts.IsGitopsSupported = updatedData.IsGitopsSupported.ValueBool()
	opts.IsHelmVMDownloadEnabled = updatedData.IsHelmvmDownloadEnabled.ValueBool()
	opts.IsIdentityServiceSupported = updatedData.IsIdentityServiceSupported.ValueBool()
	opts.IsKotsInstallEnabled = updatedData.IsKotsInstallEnabled.ValueBool()
	opts.IsSnapshotSupported = updatedData.IsSnapshotSupported.ValueBool()
	opts.IsSupportBundleUploadEnabled = updatedData.IsSupportBundleUploadEnabled.ValueBool()
	opts.Name = updatedData.Name.ValueString()
	opts.Type = updatedData.Type.ValueString()

	customerId := strings.Split(oldData.Id.ValueString(), "/")[3]

	err := r.client.UpdateCustomer(customerId, opts)

	if err != nil {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to update customer, got error: %s", err))
		return
	}

	updatedData.Id = oldData.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedData)...)
}

func (r *CustomerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CustomerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	customerId := strings.Split(data.Id.ValueString(), "/")[3]
	err := r.kotsClient.ArchiveCustomer(customerId)
	if err != nil {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to archive customer, got error: %s", err))
		return
	}
}

func (r *CustomerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
