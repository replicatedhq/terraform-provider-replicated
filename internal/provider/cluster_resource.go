// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkg/errors"
	"github.com/replicatedhq/replicated/pkg/kotsclient"
	rtypes "github.com/replicatedhq/replicated/pkg/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ClusterResource{}
var _ resource.ResourceWithImportState = &ClusterResource{}

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

// ClusterResource defines the resource implementation.
type ClusterResource struct {
	client *kotsclient.VendorV3Client
}

// ClusterResourceModel describes the resource data model.
type ClusterResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Distribution types.String `tfsdk:"distribution"`
	Version      types.String `tfsdk:"version"`
	InstanceType types.String `tfsdk:"instance_type"`
	Disk         types.Int64  `tfsdk:"disk"`
	Nodes        types.Int64  `tfsdk:"nodes"`
	TTL          types.String `tfsdk:"ttl"`
	WaitDuration types.String `tfsdk:"wait_duration"`
	Kubeconfig   types.String `tfsdk:"kubeconfig"`
}

func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Cluster resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Cluster name (defaults to random name)",
				Computed:            true,
				Optional:            true,
			},
			"distribution": schema.StringAttribute{
				MarkdownDescription: "Kubernetes distribution of the cluster to provision",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Kubernetes version to provision (format is distribution dependent)",
				Optional:            true,
				Computed:            true,
			},
			"instance_type": schema.StringAttribute{
				MarkdownDescription: "The type of instance to use (e.g. m6i.large)",
				Optional:            true,
			},
			"disk": schema.Int64Attribute{
				MarkdownDescription: "Disk Size (GiB) to request per node (default 50)",
				Optional:            true,
				Computed:            true,
			},
			"nodes": schema.Int64Attribute{
				MarkdownDescription: "Node count (default 1)",
				Optional:            true,
				Computed:            true,
			},
			"ttl": schema.StringAttribute{
				MarkdownDescription: "Cluster TTL (duration, max 48h",
				Optional:            true,
			},
			"wait_duration": schema.StringAttribute{
				MarkdownDescription: "How long to wait for the cluster to be ready",
				Optional:            true,
			},
			"kubeconfig": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kotsclient.VendorV3Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *kotsclient.VendorV3Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var waitDuration time.Duration
	if data.WaitDuration.ValueString() != "" {
		d, err := time.ParseDuration(data.WaitDuration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid wait duration", fmt.Sprintf("Unable to parse wait duration, got error: %s", err))
			return
		}
		waitDuration = d
	}

	var ttl time.Duration
	if data.TTL.ValueString() != "" {
		t, err := time.ParseDuration(data.TTL.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid ttl", fmt.Sprintf("Unable to parse ttl, got error: %s", err))
			return
		}
		ttl = t
	}

	opts := kotsclient.CreateClusterOpts{
		KubernetesDistribution: data.Distribution.ValueString(),
	}
	if data.Name.ValueString() != "" {
		opts.Name = data.Name.ValueString()
	}
	if data.Version.ValueString() != "" {
		opts.KubernetesVersion = data.Version.ValueString()
	}
	if data.InstanceType.ValueString() != "" {
		opts.InstanceType = data.InstanceType.ValueString()
	}
	if data.Disk.ValueInt64() > 0 {
		opts.DiskGiB = data.Disk.ValueInt64()
	}
	if data.Nodes.ValueInt64() > 0 {
		opts.NodeCount = int(data.Nodes.ValueInt64())
	}
	if ttl > 0 {
		opts.TTL = data.TTL.ValueString()
	}

	cl, ve, err := r.client.CreateCluster(opts)
	if err != nil {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}
	if ve != nil {
		resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("Unable to create cluster, got error: %v", ve))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// save cluster id to state
	data.Id = types.StringValue(cl.ID)
	data.Name = types.StringValue(cl.Name)
	data.Version = types.StringValue(cl.KubernetesVersion)
	data.Disk = types.Int64Value(cl.DiskGiB)
	data.Nodes = types.Int64Value(int64(cl.NodeCount))

	tflog.Trace(ctx, "created a cluster")

	// if the wait flag was provided, we poll the api until the cluster is ready, or a timeout
	if waitDuration > 0 {
		c, err := waitForCluster(r.client, cl.ID, waitDuration)
		if err != nil {
			resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
			return
		}
		if c.Status == rtypes.ClusterStatusRunning {
			k, err := r.client.GetClusterKubeconfig(c.ID)
			if err != nil {
				resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to get cluster kubeconfig, got error: %s", err))
				return
			}
			data.Kubeconfig = types.StringValue(string(k))
		}
	} else {
		data.Kubeconfig = types.StringValue("")
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cl, err := r.client.GetCluster(data.Id.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// save cluster id to state
	data.Id = types.StringValue(cl.ID)
	data.Name = types.StringValue(cl.Name)
	data.Distribution = types.StringValue(cl.KubernetesDistribution)
	data.Version = types.StringValue(cl.KubernetesVersion)
	data.Disk = types.Int64Value(cl.DiskGiB)
	data.Nodes = types.Int64Value(int64(cl.NodeCount))

	// if the cluster is running, get the kubeconfig
	if cl.Status == rtypes.ClusterStatusRunning {
		k, err := r.client.GetClusterKubeconfig(cl.ID)
		if err != nil {
			resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to get cluster kubeconfig, got error: %s", err))
			return
		}
		data.Kubeconfig = types.StringValue(string(k))
	} else {
		data.Kubeconfig = types.StringValue("")
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveCluster(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cluster, got error: %s", err))
		return
	}
}

func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForCluster(kotsRestClient *kotsclient.VendorV3Client, id string, duration time.Duration) (*rtypes.Cluster, error) {
	start := time.Now()
	for {
		cluster, err := kotsRestClient.GetCluster(id)
		if err != nil {
			return nil, errors.Wrap(err, "get cluster")
		}

		if cluster.Status == rtypes.ClusterStatusRunning {
			return cluster, nil
		} else if cluster.Status == rtypes.ClusterStatusError || cluster.Status == rtypes.ClusterStatusUpgradeError {
			return nil, errors.New("cluster failed to provision")
		} else {
			if time.Now().After(start.Add(duration)) {
				return cluster, nil
			}
		}

		time.Sleep(time.Second * 5)
	}
}
