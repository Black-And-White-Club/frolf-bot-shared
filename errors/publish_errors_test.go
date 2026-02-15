package errors

import (
	"strings"
	"testing"
)

func TestSanitizeSensitiveText(t *testing.T) {
	input := "authorization=Bearer abc token:xyz password=hunter2 https://user:pass@example.com/path"
	got := sanitizeSensitiveText(input)

	for _, wantAbsent := range []string{"abc", "xyz", "hunter2", "user:pass@"} {
		if strings.Contains(got, wantAbsent) {
			t.Fatalf("expected %q to be redacted, got %q", wantAbsent, got)
		}
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("expected redacted marker, got %q", got)
	}
}
