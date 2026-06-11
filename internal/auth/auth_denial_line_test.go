package auth

import "testing"

func TestHostDenialReason(t *testing.T) {
	tests := []struct {
		classified string
		want       string
	}{
		{"user_not_allowed", HostDenialUserNotAllowed},
		{"cli_not_enabled", HostDenialCLINotEnabled},
		{"unknown", HostDenialCLINotEnabled},
		{"user_forbidden", HostDenialUserForbidden},
		{"channel_not_allowed", HostDenialAuthDenied},
		{"channel_required", HostDenialAuthDenied},
		{"no_auth", HostDenialAuthDenied},
	}
	for _, tc := range tests {
		if got := HostDenialReason(tc.classified); got != tc.want {
			t.Fatalf("HostDenialReason(%q) = %q, want %q", tc.classified, got, tc.want)
		}
	}
}

func TestFormatAuthDenialLine(t *testing.T) {
	got := FormatAuthDenialLine("user_not_allowed")
	want := "DWS_AUTH_DENIAL reason=user_not_allowed"
	if got != want {
		t.Fatalf("FormatAuthDenialLine() = %q, want %q", got, want)
	}
}
