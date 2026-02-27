package provider

import (
	"context"
	"encoding/json"
	"testing"

	resend "github.com/resend/resend-go/v3"
)

func TestFlattenRecords_empty(t *testing.T) {
	result, diags := flattenRecords(context.Background(), []resend.Record{})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(result.Elements()) != 0 {
		t.Errorf("expected 0 elements, got %d", len(result.Elements()))
	}
}

func TestFlattenRecords_single(t *testing.T) {
	records := []resend.Record{
		{
			Record:   "SPF",
			Name:     "send",
			Type:     "TXT",
			Ttl:      "Auto",
			Status:   "verified",
			Value:    "v=spf1 include:amazonses.com ~all",
			Priority: json.Number("10"),
		},
	}

	result, diags := flattenRecords(context.Background(), records)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(result.Elements()) != 1 {
		t.Errorf("expected 1 element, got %d", len(result.Elements()))
	}
}

func TestFlattenRecords_noPriority(t *testing.T) {
	records := []resend.Record{
		{
			Record: "DKIM",
			Name:   "resend._domainkey",
			Type:   "CNAME",
			Ttl:    "Auto",
			Status: "pending",
			Value:  "some-dkim-value",
		},
	}

	result, diags := flattenRecords(context.Background(), records)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(result.Elements()) != 1 {
		t.Errorf("expected 1 element, got %d", len(result.Elements()))
	}
}

func TestFlattenRecords_multiple(t *testing.T) {
	records := []resend.Record{
		{
			Record:   "MX",
			Name:     "send",
			Type:     "MX",
			Ttl:      "Auto",
			Status:   "verified",
			Value:    "feedback-smtp.us-east-1.amazonses.com",
			Priority: json.Number("10"),
		},
		{
			Record: "SPF",
			Name:   "send",
			Type:   "TXT",
			Ttl:    "Auto",
			Status: "verified",
			Value:  "v=spf1 include:amazonses.com ~all",
		},
		{
			Record: "DKIM",
			Name:   "resend._domainkey",
			Type:   "CNAME",
			Ttl:    "Auto",
			Status: "verified",
			Value:  "some-dkim-value",
		},
	}

	result, diags := flattenRecords(context.Background(), records)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(result.Elements()) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result.Elements()))
	}
}
