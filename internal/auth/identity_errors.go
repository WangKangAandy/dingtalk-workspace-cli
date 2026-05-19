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
	"encoding/json"
	"fmt"
	"io"
)

const (
	// CodeIdentityNotAuthenticated is emitted when DWS_AUTH_IDENTITY is set but no token exists.
	CodeIdentityNotAuthenticated = "IDENTITY_NOT_AUTHENTICATED"
	// CodeIdentityMismatch is emitted when login --sender-id user does not match OAuth user.
	CodeIdentityMismatch = "IDENTITY_MISMATCH"
)

// IdentityNotAuthenticatedError indicates fail-closed: no token for the requested identity.
type IdentityNotAuthenticatedError struct {
	Identity string
}

func (e *IdentityNotAuthenticatedError) Error() string {
	return fmt.Sprintf("identity %q is not authenticated; run: dws auth login --sender-id %s --device",
		e.Identity, e.Identity)
}

// ExitCode for host/connector matching (distinct from PAT exit=4).
func (e *IdentityNotAuthenticatedError) ExitCode() int {
	return 5
}

// WriteIdentityNotAuthenticatedJSON writes a single-line machine-readable error to w.
func WriteIdentityNotAuthenticatedJSON(w io.Writer, identity string) error {
	payload := map[string]string{
		"code":     CodeIdentityNotAuthenticated,
		"identity": identity,
		"senderId": identity,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}

// WriteIdentityMismatchJSON writes IDENTITY_MISMATCH for connector handling.
func WriteIdentityMismatchJSON(w io.Writer, expected, actual string) error {
	payload := map[string]string{
		"code":     CodeIdentityMismatch,
		"expected": expected,
		"actual":   actual,
		"identity": expected,
		"senderId": expected,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}
