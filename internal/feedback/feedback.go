// Package feedback contains code for providing feedback to the user.
package feedback

import (
	"fmt"
	"os"
)

var currentPhase = 0
var priorLn = true

// Error prints an error message.
func Error(message error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
}

// Phase prints a phase indicator.
func Phase(message string) {
	total := 29
	currentPhase++
	if currentPhase > total {
		panic("incorrect total phase count")
	}

	Printf("[%d/%d] %s...\n", currentPhase, total, message)
}

// Printf prints a formatted standard message.
func Printf(format string, a ...any) {
	if !priorLn {
		fmt.Println()
		priorLn = true
	}

	fmt.Printf(format, a...)
}

// Println prints a standard message string, appending a newline.
func Println(message string) {
	if !priorLn {
		fmt.Println()
		priorLn = true
	}

	fmt.Println(message)
}

// Progress prints a progress indicator.
func Progress(current int, total int) {
	percentage := 100 * float32(current) / float32(total)
	fmt.Printf("Progress: %5.2f%%\r", percentage)

	priorLn = false
}

// Warning prints an error message.
func Warning(message error) {
	fmt.Fprintf(os.Stderr, "warning: %s\n", message)
}
