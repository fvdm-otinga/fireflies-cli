// Package errors defines the exit code taxonomy and structured error model.
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type ExitCode int

const (
	ExitSuccess      ExitCode = 0
	ExitGeneralError ExitCode = 1
	ExitUsageError   ExitCode = 2
	ExitAuthError    ExitCode = 3
	ExitRateLimit    ExitCode = 4
	ExitNotFound     ExitCode = 5
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
	fmt.Fprintf(w, "%s\n", err.Error())
	return ExitGeneralError
}

// Exit runs Handle against os.Stderr and calls os.Exit with the mapped code.
func Exit(err error) {
	os.Exit(int(Handle(os.Stderr, err)))
}
