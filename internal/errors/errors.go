// Package errors defines the exit code taxonomy and structured error model.
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type ExitCode int

// CodeKind is used by Newf to select the error category.
type CodeKind int

const (
	ExitSuccess      ExitCode = 0
	ExitGeneralError ExitCode = 1
	ExitUsageError   ExitCode = 2
	ExitAuthError    ExitCode = 3
	ExitRateLimit    ExitCode = 4
	ExitNotFound     ExitCode = 5
)

const (
	KindGeneral  CodeKind = CodeKind(ExitGeneralError)
	KindUsage    CodeKind = CodeKind(ExitUsageError)
	KindAuth     CodeKind = CodeKind(ExitAuthError)
	KindRateLimit CodeKind = CodeKind(ExitRateLimit)
	KindNotFound CodeKind = CodeKind(ExitNotFound)
)

type CLIError struct {
	Code    string   `json:"code"`
	Message string   `json:"error"`
	Exit    ExitCode `json:"-"`
}

func (e *CLIError) Error() string { return e.Message }

func Auth(msg string) *CLIError {
	return &CLIError{Code: "auth", Message: msg, Exit: ExitAuthError}
}

func Usage(msg string) *CLIError {
	return &CLIError{Code: "usage", Message: msg, Exit: ExitUsageError}
}

func NotFound(msg string) *CLIError {
	return &CLIError{Code: "not_found", Message: msg, Exit: ExitNotFound}
}

func RateLimit(msg string) *CLIError {
	return &CLIError{Code: "rate_limit", Message: msg, Exit: ExitRateLimit}
}

func General(msg string) *CLIError {
	return &CLIError{Code: "error", Message: msg, Exit: ExitGeneralError}
}

// Newf constructs a CLIError with a printf-style message.
func Newf(kind CodeKind, format string, args ...any) *CLIError {
	msg := fmt.Sprintf(format, args...)
	switch kind {
	case KindAuth:
		return Auth(msg)
	case KindUsage:
		return Usage(msg)
	case KindRateLimit:
		return RateLimit(msg)
	case KindNotFound:
		return NotFound(msg)
	default:
		return General(msg)
	}
}

// FromGraphQLError maps a GraphQL (or HTTP-level) error to a typed CLIError.
// Mapping rules:
//   - contains "Unauthorized" or "unauthorized" → code 3 (auth)
//   - contains "rate limit" or "too many" → code 4 (rate limit)
//   - contains "not found" or "Not Found" → code 5 (not found)
//   - anything else → code 1 (general)
func FromGraphQLError(err error) *CLIError {
	if err == nil {
		return nil
	}
	// Pass through CLIErrors unchanged.
	var cli *CLIError
	if errors.As(err, &cli) {
		return cli
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "unauthorized") || strings.Contains(lower, "forbidden"):
		return Auth(msg)
	case strings.Contains(lower, "rate limit") || strings.Contains(lower, "too many"):
		return RateLimit(msg)
	case strings.Contains(lower, "not found"):
		return NotFound(msg)
	default:
		return General(msg)
	}
}

// Handle writes structured error JSON to stderr (when err is a CLIError)
// or a plain message, and returns the appropriate exit code. Non-CLI errors
// map to ExitGeneralError.
func Handle(w io.Writer, err error) ExitCode {
	if err == nil {
		return ExitSuccess
	}
	var cli *CLIError
	if errors.As(err, &cli) {
		_ = json.NewEncoder(w).Encode(cli)
		return cli.Exit
	}
	_, _ = fmt.Fprintf(w, "%s\n", err.Error())
	return ExitGeneralError
}

// Exit runs Handle against os.Stderr and calls os.Exit with the mapped code.
func Exit(err error) {
	os.Exit(int(Handle(os.Stderr, err)))
}
