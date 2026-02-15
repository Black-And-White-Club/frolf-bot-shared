package observability

import (
	"strings"
	"testing"
)

func TestBuildTraceSampler_UnsetOrInvalidRateUsesSDKDefault(t *testing.T) {
	if sampler := buildTraceSampler(0); sampler != nil {
		t.Fatalf("expected nil sampler for unset rate, got %q", sampler.Description())
	}
	if sampler := buildTraceSampler(-0.5); sampler != nil {
		t.Fatalf("expected nil sampler for invalid rate, got %q", sampler.Description())
	}
}

func TestBuildTraceSampler_ExplicitRatesReturnSampler(t *testing.T) {
	if sampler := buildTraceSampler(0.25); sampler == nil {
		t.Fatal("expected non-nil sampler for explicit positive rate")
	} else if desc := sampler.Description(); !strings.Contains(desc, "TraceIDRatioBased") {
		t.Fatalf("expected ratio-based sampler, got %q", desc)
	}

	if sampler := buildTraceSampler(1); sampler == nil {
		t.Fatal("expected non-nil sampler for rate 1")
	} else if desc := sampler.Description(); !strings.Contains(desc, "AlwaysOnSampler") {
		t.Fatalf("expected always-on sampler, got %q", desc)
	}
}
