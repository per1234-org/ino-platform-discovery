// Package main implements the `main` function.
package main

//import "github.com/google/go-github/v84"
import (
	"io"
	"os"

	"github.com/per1234-org/ino-platform-discovery/internal/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	// Comment to enable logging.
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)

	rootCommand := cli.Root()
	if err := rootCommand.Execute(); err != nil {
		// Cobra handles printing error messages internally, so printing `err` here would only produce redundant output.
		os.Exit(1)
	}
}
