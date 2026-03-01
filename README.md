# terraform-provider-resend

Terraform provider for managing [Resend](https://resend.com) resources such as domains and API keys.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider)

## Usage

```hcl
terraform {
  required_providers {
    resend = {
      source = "y0n0zawa/resend"
    }
  }
}

provider "resend" {
  api_key = var.resend_api_key # Or set RESEND_API_KEY environment variable
}
```

## Authentication

Set the API key via the provider block or the `RESEND_API_KEY` environment variable:

```bash
export RESEND_API_KEY="re_..."
```

## Resources and Data Sources

### Resources

- `resend_domain` - Manages a Resend domain
- `resend_api_key` - Manages a Resend API key
- `resend_domain_verification` - Triggers domain verification

### Data Sources

- `resend_domain` - Reads a Resend domain
- `resend_api_key` - Reads a Resend API key

## Example

```hcl
# Create a domain
resource "resend_domain" "example" {
  name   = "mail.example.com"
  region = "us-east-1"
}

# Verify the domain
resource "resend_domain_verification" "example" {
  domain_id = resend_domain.example.id
}
```

## Development

### Build

```bash
go build -o terraform-provider-resend
```

### Test

```bash
go test ./...
```

### Generate documentation

```bash
go generate ./...
```

## License

[MPL-2.0](LICENSE)
