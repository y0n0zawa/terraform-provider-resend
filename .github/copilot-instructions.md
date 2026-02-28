# Copilot Instructions for terraform-provider-resend

## Project Overview

This is a Terraform provider for the [Resend](https://resend.com) email API. It manages Resend resources (domains, API keys, domain verification) via Terraform.

## Technology Stack

- **Language:** Go 1.24
- **Framework:** Terraform Plugin Framework v1.18.0 (NOT the legacy SDKv2)
- **Resend SDK:** `github.com/resend/resend-go/v3`
- **Testing:** `github.com/hashicorp/terraform-plugin-testing`

## Code Style

- All code comments must be written in English. Japanese is not allowed in comments.
- Commit messages must be written in English.
- Follow test-driven development (TDD): write tests first, confirm they fail, then implement.

## Project Structure

```
internal/provider/
  provider.go             # Provider definition (ResendProvider)
  domain_resource.go      # resend_domain resource
  api_key_resource.go     # resend_api_key resource
  domain_verification_resource.go  # resend_domain_verification resource
  domain_data_source.go   # resend_domain data source
  api_key_data_source.go  # resend_api_key data source
  errors.go               # Shared error helpers (retryOnRateLimit, isNotFoundError)
  *_test.go               # Tests
```

## Key Patterns

### API Call Retry

All Resend API calls must be wrapped with `retryOnRateLimit` (defined in `errors.go`). This generic function handles HTTP 429 rate limits and transient 500 errors with exponential backoff.

```go
result, err := retryOnRateLimit(ctx, func() (ResultType, error) {
    return r.client.Domains.GetWithContext(ctx, id)
})
```

### 404 Handling

Use `isNotFoundError(err)` to detect resources deleted outside Terraform. In `Read`, call `resp.State.RemoveResource(ctx)` when not found. In `Delete`, silently return when not found.

```go
if isNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
```

### ForceNew (RequiresReplace)

Immutable fields use `stringplanmodifier.RequiresReplace()` as a plan modifier instead of SDKv2-style `ForceNew: true`.

```go
PlanModifiers: []planmodifier.String{
    stringplanmodifier.RequiresReplace(),
},
```

### Schema Conventions

- **Computed fields** that are set once and never change use `stringplanmodifier.UseStateForUnknown()`.
- **Enum fields** use `stringvalidator.OneOf(...)` from `terraform-plugin-framework-validators`.
- **Defaults** use `stringdefault.StaticString(...)` or `booldefault.StaticBool(...)`.
- Struct tags use `tfsdk:"field_name"` for Terraform state mapping.

```go
"region": schema.StringAttribute{
    Optional: true,
    Computed: true,
    Default:  stringdefault.StaticString("us-east-1"),
    Validators: []validator.String{
        stringvalidator.OneOf("us-east-1", "eu-west-1", "sa-east-1", "ap-northeast-1"),
    },
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
},
```

### Resource Structure

Each resource follows this pattern:

1. Interface compliance: `var _ resource.Resource = &myResource{}`
2. Constructor: `func NewMyResource() resource.Resource`
3. Private struct with `client *resend.Client`
4. Model struct with `tfsdk` tags
5. Methods: `Metadata`, `Schema`, `Configure`, `Create`, `Read`, `Update`, `Delete`
6. Import support via `resource.ResourceWithImportState` and `ImportStatePassthroughID`

### Configure Pattern

Resources and data sources extract `*resend.Client` from `req.ProviderData`:

```go
func (r *myResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
```

## Testing Patterns

### Acceptance Tests

- Test function names follow `TestAcc<Resource>_<scenario>` (e.g., `TestAccDomainResource_basic`).
- Test resource names use `tf-test-` prefix for identification in sweep cleanup.
- Use `acctest.RandStringFromCharSet` for unique resource names.
- Provider factories use Protocol v6: `testAccProtoV6ProviderFactories`.
- `testAccPreCheck` ensures `RESEND_API_KEY` is set.

```go
func TestAccDomainResource_basic(t *testing.T) {
    domainName := fmt.Sprintf("tf-test-%s.example.com", acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum))
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccConfig(domainName),
                Check:  resource.ComposeAggregateTestCheckFunc(...),
            },
        },
    })
}
```

### Unit Tests

- Internal package functions are tested directly (e.g., `TestFlattenRecords_empty`).
- Unit tests live in `internal/provider/` with `_unit_test.go` suffix for non-acceptance tests within the package.

### Sweep Functions

- Sweep functions clean up leftover test resources with the `tf-test-` name prefix.
- Registered via `resource.AddTestSweepers` in `sweep_test.go`.
- `TestMain` calls `resource.TestMain(m)` to enable sweep execution.

## Error Handling

- Use `resp.Diagnostics.AddError(summary, detail)` for user-facing errors.
- Error summaries follow the pattern: `"Error <action> <resource>"`.
- Error details include the resource identifier and the original error message.
