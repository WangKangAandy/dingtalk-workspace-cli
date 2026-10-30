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
	"context"
	"strings"
	"unicode"

	authpkg "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	apperrors "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/errors"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/edition"
	"github.com/spf13/cobra"
)

// authExpectedIdentity constrains a fresh authorization, not the selection of
// an existing Profile. Never fill token identity from this caller-supplied value.
type authExpectedIdentity struct {
	corpID string
	userID string
}

func resolveExpectedAuthLoginIdentity(cmd *cobra.Command, token, profile string) (authExpectedIdentity, error) {
	if !cmd.Flags().Changed("expect-identity") {
		return authExpectedIdentity{}, nil
	}
	raw, err := cmd.Flags().GetString("expect-identity")
	if err != nil {
		return authExpectedIdentity{}, apperrors.NewInternal("failed to read --expect-identity")
	}
	corpID, userID, ok := strings.Cut(raw, ":")
	if !ok || corpID == "" || userID == "" || strings.ContainsAny(userID, ":,") ||
		strings.Contains(corpID, ",") || strings.ContainsFunc(raw, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsControl(r)
	}) {
		return authExpectedIdentity{}, apperrors.NewValidation("--expect-identity requires one exact corpId:userId (not a profile name or a list)")
	}
	if strings.TrimSpace(profile) != "" || strings.TrimSpace(token) != "" {
		return authExpectedIdentity{}, apperrors.NewValidation("--expect-identity cannot be combined with --profile or --token; use an OAuth/device login with the expected identity")
	}
	return authExpectedIdentity{corpID: corpID, userID: userID}, nil
}

func enrichExpectedAuthLoginIdentity(
	ctx context.Context,
	configDir string,
	caller edition.ToolCaller,
	data *authpkg.TokenData,
	expected authExpectedIdentity,
	hint authLoginHistoryHint,
) error {
	if err := enrichAuthLoginProfileFromContact(ctx, configDir, caller, data, hint); err != nil {
		return err
	}
	if expected.corpID == "" {
		return nil
	}
	// Both OAuth and device providers invoke this enricher before credential
	// persistence. Missing authoritative identity is a rejection, including MCP
	// logins; historical names and the requested identity cannot prove who scanned.
	if data == nil || strings.TrimSpace(data.CorpID) == "" || strings.TrimSpace(data.UserID) == "" {
		return apperrors.NewAuth("authorization did not provide a verifiable organization and user identity; credential was not saved",
			apperrors.WithReason("identity_unverified"))
	}
	if strings.TrimSpace(data.CorpID) != expected.corpID || strings.TrimSpace(data.UserID) != expected.userID {
		return apperrors.NewAuth("authorized account does not match the expected identity; credential was not saved; the requesting user must authorize with their own account",
			apperrors.WithReason("identity_mismatch"))
	}
	return nil
}
