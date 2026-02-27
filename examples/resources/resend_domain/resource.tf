resource "resend_domain" "example" {
  name   = "mail.example.com"
  region = "us-east-1"

  open_tracking  = true
  click_tracking = true
  tls            = "enforced"
}
