resource "resend_api_key" "example" {
  name       = "production"
  permission = "full_access"
}

# API key restricted to a specific domain
resource "resend_api_key" "sending_only" {
  name       = "sending-only"
  permission = "sending_access"
  domain_id  = resend_domain.example.id
}
