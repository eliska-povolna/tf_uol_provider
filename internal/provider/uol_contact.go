package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Define the contactResource struct
type contactResource struct {
	client *Client
}

// Define the data model for the resource
type contactResourceModel struct {
	Name types.String `tfsdk:"name"`
	ID   types.String `tfsdk:"id"`
}

// Implement Metadata method
func (r *contactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "uol_contact"
}

// Define the schema
func (r *contactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the contact",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "ID of the contact",
				Required:    false,
			},
		},
	}
}

// Implement Create method
func (r *contactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan contactResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactData := map[string]string{
		"name": plan.Name.ValueString(),
	}

	body, err := json.Marshal(contactData)
	if err != nil {
		resp.Diagnostics.AddError("Error encoding request body", err.Error())
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Client not initialized", "The client is nil. Ensure the client is properly initialized before making requests.")
		return
	}

	apiURL := "https://test.ucetnictvi.uol.cz/api/v1/contacts"
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		resp.Diagnostics.AddError("Error creating HTTP request", err.Error())
		return
	}

	httpResp, err := r.client.makeRequest(httpReq)
	if err != nil || httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating contact", fmt.Sprintf("API call failed with status code: %d", httpResp.StatusCode))
		return
	}

	// var response map[string]interface{}
	// json.NewDecoder(httpResp.Body).Decode(&response)
	// contactID, ok := response["id"].(string)
	// if !ok {
	// 	resp.Diagnostics.AddError("Error parsing response", response)
	// 	return
	// }

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Implement Read method
func (r *contactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state contactResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// In a real-world scenario, you'd likely fetch the contact details using the API and update the state.
	// For now, we're just assuming that the contact is always present and doesn't need to be fetched again.

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Implement Update method
func (r *contactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan contactResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactData := map[string]string{
		"name": plan.Name.ValueString(),
	}

	body, err := json.Marshal(contactData)
	if err != nil {
		resp.Diagnostics.AddError("Error encoding request body", err.Error())
		return
	}

	apiURL := fmt.Sprintf("https://test.ucetnictvi.uol.cz/api/v1/contacts/%s", plan.ID.ValueString())
	httpReq, err := http.NewRequest("PATCH", apiURL, bytes.NewBuffer(body))
	if err != nil {
		resp.Diagnostics.AddError("Error creating HTTP request", err.Error())
		return
	}

	httpResp, err := r.client.makeRequest(httpReq)
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating contact", fmt.Sprintf("API call failed with status code: %d", httpResp.StatusCode))
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Implement Delete method
func (r *contactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Error deleting contact", "The delete operation is not supported for the contact resource.")

}

// Constructor for contactResource
func NewContactResource(client *Client) *contactResource {
	return &contactResource{
		client: client,
	}
}
