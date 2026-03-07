package posm

import "testing"

func TestSanitizeAddress(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{name: "empty", in: "", out: ""},
		{name: "collapse spaces", in: "  123   Main   St  ", out: "123 Main St"},
		{name: "remove space before comma", in: "Main St , SF", out: "Main St, SF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeAddress(tt.in)
			if got != tt.out {
				t.Fatalf("SanitizeAddress(%q) = %q, want %q", tt.in, got, tt.out)
			}
		})
	}
}
