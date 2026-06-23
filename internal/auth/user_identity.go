// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	// DefaultIdentity is the legacy single-user storage slot (~/.dws root).
	DefaultIdentity = "default"

	// EnvAuthIdentity selects the active user identity (OpenClaw connector sets this).
	EnvAuthIdentity = "DWS_AUTH_IDENTITY"
)

var (
	cliIdentityMu     sync.RWMutex
	cliIdentity       string
	identitySanitizer = regexp.MustCompile(`[^a-zA-Z0-9._:@-]+`)
)

// SetCLIActiveIdentity records the --sender-id flag for the current process.
// Highest priority in ResolveActiveIdentity.
func SetCLIActiveIdentity(senderID string) {
	cliIdentityMu.Lock()
	defer cliIdentityMu.Unlock()
	cliIdentity = SanitizeIdentityID(senderID)
}

// ResolveActiveIdentity returns the identity for token storage and API auth.
// Priority: CLI --sender-id > DWS_AUTH_IDENTITY env > default.
func ResolveActiveIdentity() string {
	cliIdentityMu.RLock()
	fromCLI := cliIdentity
	cliIdentityMu.RUnlock()
	if !IsDefaultIdentity(fromCLI) {
		return fromCLI
	}
	if v := strings.TrimSpace(os.Getenv(EnvAuthIdentity)); v != "" {
		return SanitizeIdentityID(v)
	}
	return DefaultIdentity
}

// IsDefaultIdentity reports whether identity uses the legacy root config dir.
func IsDefaultIdentity(identity string) bool {
	identity = strings.TrimSpace(identity)
	return identity == "" || identity == DefaultIdentity
}

// SanitizeIdentityID normalizes an identity for use as a single path segment.
func SanitizeIdentityID(identity string) string {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return DefaultIdentity
	}
	identity = identitySanitizer.ReplaceAllString(identity, "_")
	for strings.Contains(identity, "..") {
		identity = strings.ReplaceAll(identity, "..", "_")
	}
	if identity == "" || identity == "." || identity == ".." {
		return DefaultIdentity
	}
	return identity
}

// ConfigDirForIdentity maps root ~/.dws to a per-user directory when needed.
func ConfigDirForIdentity(rootDir, identity string) string {
	rootDir = strings.TrimSpace(rootDir)
	if rootDir == "" {
		rootDir = "."
	}
	if IsDefaultIdentity(identity) {
		return rootDir
	}
	return filepath.Join(rootDir, "users", SanitizeIdentityID(identity))
}

// IsPerIdentityConfigDir reports whether configDir is under root/users/ (isolated token).
func IsPerIdentityConfigDir(rootDir, configDir string) bool {
	rootDir = filepath.Clean(strings.TrimSpace(rootDir))
	configDir = filepath.Clean(strings.TrimSpace(configDir))
	if rootDir == configDir {
		return false
	}
	usersRoot := filepath.Join(rootDir, "users")
	if configDir == usersRoot {
		return true
	}
	prefix := usersRoot + string(filepath.Separator)
	return strings.HasPrefix(configDir, prefix)
}

// IdentityIDsMatch compares connector senderId with OAuth profile id fields.
func IdentityIDsMatch(expected, actual string) bool {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)
	if expected == "" || actual == "" {
		return false
	}
	if expected == actual {
		return true
	}
	// corpId:senderId vs bare senderId
	if i := strings.LastIndex(expected, ":"); i >= 0 {
		if strings.TrimSpace(expected[i+1:]) == actual {
			return true
		}
	}
	if i := strings.LastIndex(actual, ":"); i >= 0 {
		if strings.TrimSpace(actual[i+1:]) == expected {
			return true
		}
	}
	return false
}
