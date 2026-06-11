// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package errors

import (
	"context"
	stderrors "errors"
	"strings"
)

// ExitCodeLoginTimeout is the process exit code for device/OAuth login timeouts.
const ExitCodeLoginTimeout = 5

// LoginTimeoutError is returned when login exceeds its deadline or authorization code expires.
type LoginTimeoutError struct {
	Message string
	Cause   error
}

func (e *LoginTimeoutError) Error() string {
	if e == nil {
		return "login timed out"
	}
	if e.Message != "" {
		return e.Message
	}
	return "login timed out"
}

func (e *LoginTimeoutError) Unwrap() error { return e.Cause }

// ExitCode implements ExitCoder for connector / host matching.
func (e *LoginTimeoutError) ExitCode() int { return ExitCodeLoginTimeout }

// NewLoginTimeout wraps a timeout cause as LoginTimeoutError (exit 5).
func NewLoginTimeout(cause error) error {
	msg := "login timed out"
	if cause != nil {
		msg = cause.Error()
	}
	return &LoginTimeoutError{Message: msg, Cause: cause}
}

// IsLoginTimeout reports whether err is a login timeout (deadline or known timeout text).
func IsLoginTimeout(err error) bool {
	if err == nil {
		return false
	}
	var lt *LoginTimeoutError
	if stderrors.As(err, &lt) {
		return true
	}
	if stderrors.Is(err, context.DeadlineExceeded) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "操作超时") ||
		strings.Contains(msg, "context deadline exceeded") ||
		strings.Contains(msg, "authorization code expired") ||
		strings.Contains(msg, "授权码已过期")
}
