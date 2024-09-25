// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure uolProvider satisfies various provider interfaces.
var _ provider.Provider = &uolProvider{}

// uolProvider defines the provider implementation.
type uolProvider struct {
	client *Client
}

// uolProviderModel describes the provider data model.
type uolProviderModel struct {
	Email types.String `tfsdk:"email"`
	Token types.String `tfsdk:"token"`
}

func (p *uolProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uol"
}

func (p *uolProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "Email for API authentication",
				Required:    true,
			},
			"token": schema.StringAttribute{
				Description: "API token for authentication",
				Required:    true,
			},
		},
	}
}

func (p *uolProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config uolProviderModel

	// Retrieve provider configuration
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize email and token
	var email, token string

	// Check if the email is provided and valid
	if !config.Email.IsNull() {
		email = config.Email.ValueString()
	} else {
		resp.Diagnostics.AddError("Missing Email", "The provider configuration is missing the 'email' attribute.")
		return
	}

	// Check if the token is provided and valid
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	} else {
		resp.Diagnostics.AddError("Missing Token", "The provider configuration is missing the 'token' attribute.")
		return
	}
	// Set up the API client
	client := &Client{
		Email:      email,
		Token:      token,
		HttpClient: &http.Client{},
	}
	resp.DataSourceData = client
	resp.ResourceData = client
	p.client = client
}

func (p *uolProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return NewContactResource(p.client)
		},
	}
}

func (p *uolProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *uolProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

// Implement the provider
func New() provider.Provider {
	return &uolProvider{}
}
