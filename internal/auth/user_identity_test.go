// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeIdentityID(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", DefaultIdentity},
		{"user-001", "user-001"},
		{"corp:staff123", "corp:staff123"},
		{"../../etc", "____etc"},
	}
	for _, tc := range tests {
		if got := SanitizeIdentityID(tc.in); got != tc.want {
			t.Fatalf("SanitizeIdentityID(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestConfigDirForIdentity(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".dws")
	defaultDir := ConfigDirForIdentity(root, DefaultIdentity)
	if defaultDir != root {
		t.Fatalf("default dir = %q, want %q", defaultDir, root)
	}
	userDir := ConfigDirForIdentity(root, "userA")
	want := filepath.Join(root, "users", "userA")
	if userDir != want {
		t.Fatalf("user dir = %q, want %q", userDir, want)
	}
}

func TestResolveActiveIdentityPriority(t *testing.T) {
	t.Setenv(EnvAuthIdentity, "from-env")
	SetCLIActiveIdentity("from-cli")
	if got := ResolveActiveIdentity(); got != "from-cli" {
		t.Fatalf("ResolveActiveIdentity() = %q, want from-cli", got)
	}
	SetCLIActiveIdentity("")
	if got := ResolveActiveIdentity(); got != "from-env" {
		t.Fatalf("ResolveActiveIdentity() = %q, want from-env", got)
	}
	_ = os.Unsetenv(EnvAuthIdentity)
	if got := ResolveActiveIdentity(); got != DefaultIdentity {
		t.Fatalf("ResolveActiveIdentity() = %q, want default", got)
	}
}

func TestIdentityIDsMatch(t *testing.T) {
	if !IdentityIDsMatch("staff1", "staff1") {
		t.Fatal("expected exact match")
	}
	if !IdentityIDsMatch("corp:staff1", "staff1") {
		t.Fatal("expected corp-prefixed match")
	}
	if IdentityIDsMatch("staff1", "staff2") {
		t.Fatal("expected mismatch")
	}
}
