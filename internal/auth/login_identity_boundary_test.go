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
	"errors"
	"io"
	"testing"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/keychain"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/testseam"
)

// Exercise the actual provider persistence boundaries, not a fake login
// implementation: rejecting identity enrichment must precede every save.
func TestCrossPlatformCoverageLoginIdentityRejectionPrecedesPersistence(t *testing.T) {
	for _, flow := range []string{"oauth", "device"} {
		t.Run(flow, func(t *testing.T) {
			t.Setenv(keychain.TestNamespaceEnv, t.TempDir())
			cleanupKeychain(t)
			dir := t.TempDir()
			t.Setenv("DWS_CONFIG_DIR", dir)
			rejected := errors.New("identity rejected")
			incoming := &TokenData{AccessToken: "new-token", CorpID: "ding_corp", UserID: "bob"}
			enrich := func(context.Context, *TokenData) error { return rejected }
			writes := 0
			save := func(string, *TokenData) error { writes++; return nil }
			testseam.Swap(t, &oauthSaveToken, save)
			testseam.Swap(t, &deviceSaveToken, save)
			var err error
			if flow == "oauth" {
				p := &OAuthProvider{configDir: dir, IdentityEnricher: enrich}
				err = p.persistLoginToken(context.Background(), incoming)
			} else {
				testseam.Swap(t, &deviceRequestCode, func(*DeviceFlowProvider, context.Context) (*DeviceAuthResponse, error) {
					return &DeviceAuthResponse{DeviceCode: "dc", Interval: 1}, nil
				})
				testseam.Swap(t, &deviceWaitAuth, func(*DeviceFlowProvider, context.Context, *DeviceAuthResponse) (*DeviceTokenResponse, error) {
					return &DeviceTokenResponse{AuthCode: "code"}, nil
				})
				testseam.Swap(t, &deviceExchangeCode, func(*OAuthProvider, context.Context, string) (*TokenData, error) { return incoming, nil })
				testseam.Swap(t, &deviceCheckCLIAuth, func(*OAuthProvider, context.Context, string) (*CLIAuthStatus, error) {
					return &CLIAuthStatus{Success: true, Result: &CLIAuthResult{CLIAuthEnabled: true}}, nil
				})
				p := &DeviceFlowProvider{configDir: dir, clientID: "client", IdentityEnricher: enrich, Output: io.Discard, NoBrowser: true}
				_, err = p.loginOnce(context.Background(), 1)
			}
			if !errors.Is(err, rejected) || writes != 0 {
				t.Fatalf("identity rejection: err=%v saves=%d", err, writes)
			}
		})
	}
}
