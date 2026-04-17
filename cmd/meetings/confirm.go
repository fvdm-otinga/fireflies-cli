package meetings

import (
	"os"
)

// isTerminal reports whether f is connected to a terminal (TTY).
// Delegates to a stat-based check (ModeCharDevice).
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
