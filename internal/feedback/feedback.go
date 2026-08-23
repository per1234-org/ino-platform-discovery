// Package feedback contains code for providing feedback to the user.
package feedback

import (
	"fmt"
	"os"
)

// Error prints an error message.
func Error(message error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
}

// Printf prints a formatted standard message.
func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}

// Println prints a standard message string, appending a newline.
func Println(message string) {
	fmt.Println(message)
}

// Warning prints an error message.
func Warning(message error) {
	fmt.Fprintf(os.Stderr, "warning: %s\n", message)
}
