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
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	authpkg "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	apperrors "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/errors"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/keychain"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/testseam"
)

func TestCrossPlatformCoverageExpectedLoginIdentityValidation(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"empty", []string{"--expect-identity="}},
		{"name", []string{"--expect-identity=alice"}},
		{"missing-user", []string{"--expect-identity=ding_corp:"}},
		{"missing-corp", []string{"--expect-identity=:alice"}},
		{"multiple", []string{"--expect-identity=ding_corp:alice,ding_corp:bob"}},
		{"extra-colon", []string{"--expect-identity=ding_corp:alice:bob"}},
		{"whitespace", []string{"--expect-identity=ding_corp: alice"}},
		{"manual-token", []string{"--expect-identity=ding_corp:alice", "--token=secret"}},
		{"profile", []string{"--expect-identity=ding_corp:alice", "--profile=ding_corp:alice"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(keychain.TestNamespaceEnv, t.TempDir())
			setupAuthLogoutProfiles(t)
			testseam.Swap(t, &authOAuthLogin, func(*authpkg.OAuthProvider, context.Context, bool) (*authpkg.TokenData, error) {
				t.Fatal("invalid expectation started authorization")
				return nil, nil
			})
			root := NewRootCommand()
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)
			root.SetArgs(append([]string{"auth", "login", "--format=json"}, tc.args...))
			err := root.Execute()
			var classified *apperrors.Error
			if !errors.As(err, &classified) || classified.Category != apperrors.CategoryValidation {
				t.Fatalf("error = %v, want validation; output=%s", err, out.String())
			}
		})
	}
}

func TestCrossPlatformCoverageExpectedLoginIdentityOAuthAndDevice(t *testing.T) {
	for _, flow := range []string{"oauth", "device"} {
		for _, tc := range []struct {
			name, corp, user, reason string
			guard                    bool
		}{
			{"first-login", "ding_corp", "alice", "", true},
			{"wrong-user", "ding_corp", "bob", "identity_mismatch", true},
			{"wrong-corp", "ding_other", "alice", "identity_mismatch", true},
			{"unknown-user", "ding_corp", "", "identity_unverified", true},
			{"unknown-corp", "", "alice", "identity_unverified", true},
			{"legacy-unresolved", "ding_corp", "", "", false},
		} {
			t.Run(flow+"/"+tc.name, func(t *testing.T) {
				t.Setenv(keychain.TestNamespaceEnv, t.TempDir())
				setupAuthLogoutProfiles(t)
				testseam.Swap(t, &authLoginInteractiveTerminal, func() bool { return false })
				checked := false
				login := func(ctx context.Context, enrich func(context.Context, *authpkg.TokenData) error) (*authpkg.TokenData, error) {
					checked = true
					if enrich == nil {
						t.Fatal("missing pre-persistence identity hook")
					}
					data := &authpkg.TokenData{AccessToken: "must-not-leak", CorpID: tc.corp, UserID: tc.user}
					if err := enrich(ctx, data); err != nil {
						return nil, fmt.Errorf("provider persistence boundary: %w", err)
					}
					return data, nil
				}
				testseam.Swap(t, &authOAuthLogin, func(p *authpkg.OAuthProvider, ctx context.Context, _ bool) (*authpkg.TokenData, error) {
					if tc.guard && p.TargetCorpID != "ding_corp" {
						t.Fatalf("target corp = %q", p.TargetCorpID)
					}
					return login(ctx, p.IdentityEnricher)
				})
				testseam.Swap(t, &authDeviceLogin, func(p *authpkg.DeviceFlowProvider, ctx context.Context) (*authpkg.TokenData, error) {
					return login(ctx, p.IdentityEnricher)
				})
				root := NewRootCommand()
				var out bytes.Buffer
				root.SetOut(&out)
				root.SetErr(&out)
				args := []string{"auth", "login", "--format=json", "--no-browser"}
				if flow == "device" {
					args = append(args, "--device")
				}
				if tc.guard {
					args = append(args, "--expect-identity=ding_corp:alice")
				}
				root.SetArgs(args)
				err := root.Execute()
				if !checked {
					t.Fatal("authorization hook not reached")
				}
				if tc.reason == "" {
					if err != nil {
						t.Fatalf("login failed: %v; output=%s", err, out.String())
					}
				} else {
					var classified *apperrors.Error
					if !errors.As(err, &classified) || classified.Category != apperrors.CategoryAuth || classified.Reason != tc.reason {
						t.Fatalf("error = %v, want auth/%s", err, tc.reason)
					}
					if classified.ExitCode() != 2 {
						t.Fatalf("exit=%d, want 2", classified.ExitCode())
					}
					if strings.Contains(err.Error(), "must-not-leak") {
						t.Fatal("error disclosed credential")
					}
				}
			})
		}
	}
}

func TestCrossPlatformCoverageExpectedLoginIdentityErrorScope(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		for _, reason := range []string{"identity_mismatch", "identity_unverified", "other"} {
			for _, category := range []string{"auth", "api"} {
				t.Run(fmt.Sprintf("enabled=%t/%s/%s", enabled, category, reason), func(t *testing.T) {
					expected := authExpectedIdentity{}
					if enabled {
						expected = authExpectedIdentity{corpID: "corp", userID: "user"}
					}
					classified := apperrors.NewAuth("failure", apperrors.WithReason(reason))
					if category == "api" {
						classified = apperrors.NewAPI("failure", apperrors.WithReason(reason))
					}
					wrapped := fmt.Errorf("provider: %w", classified)
					got := expected.flowError("login failed", wrapped)
					if enabled && category == "auth" && reason != "other" {
						if got != wrapped {
							t.Fatalf("guard error not preserved: %v", got)
						}
						return
					}
					want := authLoginFlowError("login failed", wrapped)
					var actual, baseline *apperrors.Error
					if !errors.As(got, &actual) || !errors.As(want, &baseline) ||
						actual.Category != baseline.Category || actual.Reason != baseline.Reason ||
						actual.Error() != baseline.Error() || actual.ExitCode() != baseline.ExitCode() {
						t.Fatalf("non-guard error differs from upstream: got=%v want=%v", got, want)
					}
				})
			}
		}
	}
}

func TestCrossPlatformCoverageExpectedLoginIdentityDoesNotTrustHistory(t *testing.T) {
	t.Setenv(keychain.TestNamespaceEnv, t.TempDir())
	configDir := setupAuthLogoutProfiles(t, authLogoutTestToken("ding_corp"))
	data := &authpkg.TokenData{AccessToken: "fresh", CorpID: "ding_corp"}
	err := enrichExpectedAuthLoginIdentity(context.Background(), configDir, nil, data,
		authExpectedIdentity{corpID: "ding_corp", userID: "user-ding_corp"},
		authLoginHistoryHint{Selector: "ding_corp:user-ding_corp", Explicit: true})
	var classified *apperrors.Error
	if !errors.As(err, &classified) || classified.Reason != "identity_unverified" || data.UserID != "" {
		t.Fatalf("historical identity trusted: user=%q error=%v", data.UserID, err)
	}
}
