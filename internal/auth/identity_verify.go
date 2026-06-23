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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/apiclient"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/config"
)

// VerifyLoginIdentity ensures the OAuth token belongs to expected senderId before persisting.
// Skipped for default (non --sender-id) login.
func VerifyLoginIdentity(ctx context.Context, expected string, tokenData *TokenData) error {
	if IsDefaultIdentity(expected) || tokenData == nil {
		return nil
	}
	// MCP OAuth tokens are for mcp.dingtalk.com (x-user-access-token), not DingTalk OpenAPI
	// userAccessTokens — see `dws api` long help ("MCP 默认凭证...不支持 raw API"). Per-sender
	// mismatch checks via GET /v1.0/contact/users/me are only possible in direct OAuth mode.
	if IsClientIDFromMCP() {
		return nil
	}
	actual, err := FetchAuthenticatedStaffID(ctx, tokenData.AccessToken)
	if err != nil {
		return fmt.Errorf("verify login identity: %w", err)
	}
	if !IdentityIDsMatch(expected, actual) {
		return &IdentityMismatchError{Expected: expected, Actual: actual}
	}
	return nil
}

// FetchAuthenticatedStaffID calls GET /v1.0/contact/users/me with the access token.
func FetchAuthenticatedStaffID(ctx context.Context, accessToken string) (string, error) {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return "", fmt.Errorf("empty access token")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, UserInfoURL, nil)
	if err != nil {
		return "", err
	}
	if err := apiclient.ValidateTargetHost(UserInfoURL); err != nil {
		return "", err
	}
	// api.dingtalk.com OpenAPI (direct userAccessToken): x-acs-dingtalk-access-token.
	req.Header.Set(apiclient.AuthHeader, accessToken)

	client := &http.Client{Timeout: config.OAuthTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("user info HTTP %d: %s", resp.StatusCode, truncateIdentityBody(body))
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("parse user info: %w", err)
	}
	for _, key := range []string{"staffId", "userId", "userid", "unionId"} {
		if raw, ok := payload[key]; ok {
			var s string
			if err := json.Unmarshal(raw, &s); err == nil && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s), nil
			}
		}
	}
	return "", fmt.Errorf("user info missing staffId/userId")
}

func truncateIdentityBody(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}

// IdentityMismatchError is returned when --sender-id does not match the OAuth user.
type IdentityMismatchError struct {
	Expected string
	Actual   string
}

func (e *IdentityMismatchError) Error() string {
	return fmt.Sprintf(
		"identity mismatch: expected senderId %q but authenticated user is %q; use your own account to scan",
		e.Expected,
		e.Actual,
	)
}
