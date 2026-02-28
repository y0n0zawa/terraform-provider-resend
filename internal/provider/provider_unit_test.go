package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestNew(t *testing.T) {
	factory := New("1.0.0")
	p := factory()
	if p == nil {
		t.Error("expected non-nil provider")
	}
}

func TestProvider_Metadata(t *testing.T) {
	p := &ResendProvider{version: "1.0.0"}
	req := provider.MetadataRequest{}
	resp := provider.MetadataResponse{}
	p.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "resend" {
		t.Errorf("expected TypeName 'resend', got %q", resp.TypeName)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected Version '1.0.0', got %q", resp.Version)
	}
}

func TestProvider_Schema(t *testing.T) {
	p := &ResendProvider{}
	req := provider.SchemaRequest{}
	resp := provider.SchemaResponse{}
	p.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["api_key"] == nil {
		t.Error("expected api_key attribute in schema")
	}
}

func providerConfig(apiKey *string) tfsdk.Config {
	p := &ResendProvider{}
	var schemaResp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &schemaResp)
	objType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"api_key": tftypes.String,
		},
	}
	var apiKeyVal tftypes.Value
	if apiKey != nil {
		apiKeyVal = tftypes.NewValue(tftypes.String, *apiKey)
	} else {
		apiKeyVal = tftypes.NewValue(tftypes.String, nil)
	}
	return tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
			"api_key": apiKeyVal,
		}),
	}
}

func TestProvider_Configure_withEnvVar(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	p := &ResendProvider{}
	req := provider.ConfigureRequest{Config: providerConfig(nil)}
	resp := provider.ConfigureResponse{}
	p.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	if resp.ResourceData == nil {
		t.Error("expected ResourceData to be set")
	}
	if resp.DataSourceData == nil {
		t.Error("expected DataSourceData to be set")
	}
}

func TestProvider_Configure_withConfig(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "")
	key := "config-key"
	p := &ResendProvider{}
	req := provider.ConfigureRequest{Config: providerConfig(&key)}
	resp := provider.ConfigureResponse{}
	p.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	if resp.ResourceData == nil {
		t.Error("expected ResourceData to be set")
	}
}

func TestProvider_Configure_configOverridesEnv(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "env-key")
	key := "config-key"
	p := &ResendProvider{}
	req := provider.ConfigureRequest{Config: providerConfig(&key)}
	resp := provider.ConfigureResponse{}
	p.Configure(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
}

func TestProvider_Configure_missingKey(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "")
	p := &ResendProvider{}
	req := provider.ConfigureRequest{Config: providerConfig(nil)}
	resp := provider.ConfigureResponse{}
	p.Configure(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected error for missing API key")
	}
}

func TestProvider_Resources(t *testing.T) {
	p := &ResendProvider{}
	resources := p.Resources(context.Background())
	if len(resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(resources))
	}
}

func TestProvider_DataSources(t *testing.T) {
	p := &ResendProvider{}
	dataSources := p.DataSources(context.Background())
	if len(dataSources) != 2 {
		t.Errorf("expected 2 data sources, got %d", len(dataSources))
	}
}
