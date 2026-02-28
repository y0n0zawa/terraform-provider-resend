data "resend_api_key" "example" {
  id = "key_123456"
}

output "api_key_name" {
  value = data.resend_api_key.example.name
}
