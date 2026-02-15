package observability

import "testing"

func TestRedactLogValue(t *testing.T) {
	if got := redactLogValue("interaction_token", "secret-token"); got != "[REDACTED]" {
		t.Fatalf("expected redacted token, got %q", got)
	}
	if got := redactLogValue("request_id", "abc123"); got != "abc123" {
		t.Fatalf("expected non-sensitive value unchanged, got %q", got)
	}
}
