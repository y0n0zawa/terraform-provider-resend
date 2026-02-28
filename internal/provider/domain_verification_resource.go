package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	resend "github.com/resend/resend-go/v3"
)

var _ resource.Resource = &domainVerificationResource{}

// NewDomainVerificationResource returns a new domain verification resource.
func NewDomainVerificationResource() resource.Resource {
	return &domainVerificationResource{}
}

type domainVerificationResource struct {
	client *resend.Client
}

type domainVerificationResourceModel struct {
	DomainID types.String `tfsdk:"domain_id"`
	Status   types.String `tfsdk:"status"`
}

func (r *domainVerificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_verification"
}

func (r *domainVerificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Triggers verification for a Resend domain. Creating this resource sends a verification request to Resend. The domain's DNS records must be configured before verification can succeed.",
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the domain to verify.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The verification status of the domain after triggering verification.",
				Computed:            true,
			},
		},
	}
}

func (r *domainVerificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*resend.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *resend.Client")
		return
	}
	r.client = client
}

func (r *domainVerificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainVerificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID := plan.DomainID.ValueString()

	// Trigger verification
	_, err := retryOnRateLimit(ctx, func() (bool, error) {
		return r.client.Domains.VerifyWithContext(ctx, domainID)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error verifying domain",
			"Could not trigger verification for domain ID "+domainID+": "+err.Error(),
		)
		return
	}

	// Read domain to get current status
	domain, err := retryOnRateLimit(ctx, func() (resend.Domain, error) {
		return r.client.Domains.GetWithContext(ctx, domainID)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain after verification",
			"Verification was triggered for domain "+domainID+" but could not read status: "+err.Error(),
		)
		return
	}

	plan.Status = types.StringValue(domain.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainVerificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainVerificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := retryOnRateLimit(ctx, func() (resend.Domain, error) {
		return r.client.Domains.GetWithContext(ctx, state.DomainID.ValueString())
	})
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading domain",
			"Could not read domain ID "+state.DomainID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Status = types.StringValue(domain.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *domainVerificationResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// domain_id is ForceNew; Update should never be called.
}

func (r *domainVerificationResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Verification cannot be undone; no-op.
}
