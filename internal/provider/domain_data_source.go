package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	resend "github.com/resend/resend-go/v3"
)

var _ datasource.DataSource = &domainDataSource{}

// NewDomainDataSource returns a new domain data source.
func NewDomainDataSource() datasource.DataSource {
	return &domainDataSource{}
}

type domainDataSource struct {
	domains resend.DomainsSvc
}

type domainDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Region    types.String `tfsdk:"region"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	Records   types.List   `tfsdk:"records"`
}

// populateFromDomain sets the common fields from a Resend API domain response.
func (m *domainDataSourceModel) populateFromDomain(ctx context.Context, domain resend.Domain) diag.Diagnostics {
	m.ID = types.StringValue(domain.Id)
	m.Name = types.StringValue(domain.Name)
	m.Status = types.StringValue(domain.Status)
	m.Region = types.StringValue(domain.Region)
	m.CreatedAt = types.StringValue(domain.CreatedAt)

	records, diags := flattenRecords(ctx, domain.Records)
	if !diags.HasError() {
		m.Records = records
	}
	return diags
}

func (d *domainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *domainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to read a Resend domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The domain ID.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The domain name.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region where the domain is hosted.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The verification status of the domain.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the domain was created.",
				Computed:            true,
			},
			"records": schema.ListNestedAttribute{
				MarkdownDescription: "DNS records required for domain verification.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"record": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"ttl": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"value": schema.StringAttribute{
							Computed: true,
						},
						"priority": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *domainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.domains = configureDomainsSvc(req.ProviderData, &resp.Diagnostics)
}

func (d *domainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data domainDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := retryOnRateLimit(ctx, func() (resend.Domain, error) {
		return d.domains.GetWithContext(ctx, data.ID.ValueString())
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain",
			"Could not read domain ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.populateFromDomain(ctx, domain)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
