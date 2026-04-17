// Package confirm provides interactive confirmation prompts for destructive operations.
package confirm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

// IsTerminal reports whether f is connected to a terminal (TTY).
func IsTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// Require prompts the user to type "YES" to confirm a destructive operation
// described by msg. If yes is true the prompt is skipped. If stdin is not a
// TTY and yes is false, a usage error (exit 2) is returned.
//
// prompt is written to stderr; input is read from in (typically os.Stdin).
func Require(yes bool, in io.Reader, msg string) error {
	if yes {
		return nil
	}

	// Check if in is a terminal (os.Stdin specifically).
	if f, ok := in.(*os.File); ok && !IsTerminal(f) {
		return ferr.Usage("stdin is not a TTY and --yes is not set; cannot confirm destructive operation")
	}

	fmt.Fprintf(os.Stderr, "%s\nType YES to confirm: ", msg)
	scanner := bufio.NewScanner(in)
	scanner.Scan()
	if strings.TrimSpace(scanner.Text()) != "YES" {
		return ferr.Usage("operation cancelled")
	}
	return nil
}
