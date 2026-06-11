// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package auth

import (
	"fmt"
	"io"
)

// Host-facing denial reasons emitted on the DWS_AUTH_DENIAL stderr contract line.
// See dingtalk-openclaw-connector docs/DWS_AUTH_PHASE1_5_REFACTOR.md.
const (
	HostDenialUserNotAllowed  = "user_not_allowed"
	HostDenialCLINotEnabled   = "cli_not_enabled"
	HostDenialUserForbidden   = "user_forbidden"
	HostDenialAuthDenied      = "auth_denied"
)

// HostDenialReason maps classifyDenialReason output to the stable host contract.
func HostDenialReason(classified string) string {
	switch classified {
	case "user_not_allowed":
		return HostDenialUserNotAllowed
	case "cli_not_enabled", "unknown":
		return HostDenialCLINotEnabled
	case "user_forbidden":
		return HostDenialUserForbidden
	default:
		return HostDenialAuthDenied
	}
}

// FormatAuthDenialLine returns the machine-readable stderr line for hosts/connectors.
func FormatAuthDenialLine(classifiedReason string) string {
	return fmt.Sprintf("DWS_AUTH_DENIAL reason=%s", HostDenialReason(classifiedReason))
}

// WriteAuthDenialLine prints the contract line to w (typically login stderr).
func WriteAuthDenialLine(w io.Writer, classifiedReason string) {
	_, _ = fmt.Fprintln(w, FormatAuthDenialLine(classifiedReason))
}
