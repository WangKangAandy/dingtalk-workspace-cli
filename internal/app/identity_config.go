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

package app

import (
	"errors"
	"io"
	"os"

	authpkg "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
)

// resolvedConfigDir returns the token storage directory for the active identity.
func resolvedConfigDir() string {
	return authpkg.ConfigDirForIdentity(defaultConfigDir(), authpkg.ResolveActiveIdentity())
}

// syncActiveIdentityFromFlags applies --sender-id and keeps DWS_AUTH_IDENTITY aligned for cache partition.
func syncActiveIdentityFromFlags(flags *GlobalFlags) {
	senderID := ""
	if flags != nil {
		senderID = flags.SenderID
	}
	authpkg.SetCLIActiveIdentity(senderID)
	identity := authpkg.ResolveActiveIdentity()
	if authpkg.IsDefaultIdentity(identity) {
		_ = os.Unsetenv(authpkg.EnvAuthIdentity)
		return
	}
	_ = os.Setenv(authpkg.EnvAuthIdentity, identity)
}

// noCredentialsErrorForActiveIdentity returns fail-closed errors when identity is explicit.
func noCredentialsErrorForActiveIdentity() error {
	identity := authpkg.ResolveActiveIdentity()
	if authpkg.IsDefaultIdentity(identity) {
		return noCredentialsError()
	}
	return &authpkg.IdentityNotAuthenticatedError{Identity: identity}
}

// writeIdentityAuthError emits machine-readable JSON for connector when appropriate.
func writeIdentityAuthError(w io.Writer, err error) bool {
	var notAuth *authpkg.IdentityNotAuthenticatedError
	if errors.As(err, &notAuth) {
		_ = authpkg.WriteIdentityNotAuthenticatedJSON(w, notAuth.Identity)
		return true
	}
	var mismatch *authpkg.IdentityMismatchError
	if errors.As(err, &mismatch) {
		_ = authpkg.WriteIdentityMismatchJSON(w, mismatch.Expected, mismatch.Actual)
		return true
	}
	return false
}
