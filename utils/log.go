package utils

import (
	"fmt"
	"os"
)

// PrintDebug prints debug messages.
func PrintDebug(msg string, debug bool) {
	if debug {
		fmt.Fprintln(os.Stderr, "DEBUG: "+msg)
	}
}

// PrintWarning prints warning messages.
func PrintWarning(msg string, debug bool) {
	if debug {
		fmt.Fprintln(os.Stderr, "Warning: "+msg)
	}
}

// PrintError prints error messages.
func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, "Error: "+msg)
}