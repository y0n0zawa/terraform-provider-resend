package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resend "github.com/resend/resend-go/v3"
)

// configureClient extracts the Resend client from provider data.
func configureClient(providerData any, diagnostics *diag.Diagnostics) *resend.Client {
	if providerData == nil {
		return nil
	}
	client, ok := providerData.(*resend.Client)
	if !ok {
		diagnostics.AddError("Unexpected Configure Type", "Expected *resend.Client")
		return nil
	}
	return client
}
