resource "resend_domain" "example" {
  name = "example.com"
}

resource "resend_domain_verification" "example" {
  domain_id = resend_domain.example.id
}
