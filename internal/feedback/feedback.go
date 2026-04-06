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

// Warning prints an error message.
func Warning(message error) {
	fmt.Fprintf(os.Stderr, "warning: %s\n", message)
}
