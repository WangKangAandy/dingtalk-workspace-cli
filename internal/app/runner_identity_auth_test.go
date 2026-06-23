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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authpkg "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	apperrors "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/errors"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/executor"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/transport"
)

func TestRuntimeRunner_ExplicitIdentity_NoToken_ReturnsIdentityNotAuthenticated(t *testing.T) {
	setupRuntimeCommandTest(t)
	t.Setenv(authpkg.EnvAuthIdentity, "user-without-token")
	authpkg.SetCLIActiveIdentity("")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("HTTP must not be called when identity has no token")
	}))
	t.Cleanup(server.Close)

	runner := &runtimeRunner{
		transport: transport.NewClient(nil),
		fallback:  executor.EchoRunner{},
	}

	_, err := runner.executeInvocation(context.Background(), server.URL, executor.Invocation{
		CanonicalProduct: "doc",
		Tool:             "search_documents",
		Params:           map[string]any{"keyword": "test"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var identityErr *authpkg.IdentityNotAuthenticatedError
	if !errors.As(err, &identityErr) {
		t.Fatalf("expected IdentityNotAuthenticatedError, got %T: %v", err, err)
	}
	if identityErr.Identity != "user-without-token" {
		t.Fatalf("identity = %q, want user-without-token", identityErr.Identity)
	}
	if identityErr.ExitCode() != 5 {
		t.Fatalf("ExitCode() = %d, want 5", identityErr.ExitCode())
	}
}

func TestRuntimeRunner_DefaultIdentity_NoToken_ReturnsGenericAuthError(t *testing.T) {
	setupRuntimeCommandTest(t)
	t.Setenv(authpkg.EnvAuthIdentity, "")
	authpkg.SetCLIActiveIdentity("")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("HTTP must not be called when default identity has no token")
	}))
	t.Cleanup(server.Close)

	runner := &runtimeRunner{
		transport: transport.NewClient(nil),
		fallback:  executor.EchoRunner{},
	}

	_, err := runner.executeInvocation(context.Background(), server.URL, executor.Invocation{
		CanonicalProduct: "doc",
		Tool:             "search_documents",
		Params:           map[string]any{"keyword": "test"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var identityErr *authpkg.IdentityNotAuthenticatedError
	if errors.As(err, &identityErr) {
		t.Fatalf("default identity should not return IdentityNotAuthenticatedError, got %v", identityErr)
	}
	var authErr *apperrors.Error
	if !errors.As(err, &authErr) || authErr.Category != apperrors.CategoryAuth {
		t.Fatalf("expected auth-category error, got %T: %v", err, err)
	}
	if authErr.Reason != "not_authenticated" {
		t.Fatalf("reason = %q, want not_authenticated", authErr.Reason)
	}
}
