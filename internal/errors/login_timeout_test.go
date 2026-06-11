package errors

import (
	"context"
	"errors"
	"testing"
)

func TestLoginTimeoutError_ExitCode(t *testing.T) {
	err := NewLoginTimeout(context.DeadlineExceeded)
	var lt *LoginTimeoutError
	if !errors.As(err, &lt) {
		t.Fatal("expected LoginTimeoutError")
	}
	if lt.ExitCode() != ExitCodeLoginTimeout {
		t.Fatalf("ExitCode() = %d, want %d", lt.ExitCode(), ExitCodeLoginTimeout)
	}
	if ExitCode(err) != ExitCodeLoginTimeout {
		t.Fatalf("ExitCode(err) = %d, want %d", ExitCode(err), ExitCodeLoginTimeout)
	}
}

func TestIsLoginTimeout(t *testing.T) {
	if !IsLoginTimeout(context.DeadlineExceeded) {
		t.Fatal("expected deadline exceeded to be login timeout")
	}
	if !IsLoginTimeout(NewLoginTimeout(errors.New("操作超时，请重新登录"))) {
		t.Fatal("expected Chinese timeout hint")
	}
	if IsLoginTimeout(NewAuth("user not allowed")) {
		t.Fatal("auth error should not be login timeout")
	}
}
