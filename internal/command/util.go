package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Check if err != nil. If so, print the error, command usage (if printUsage is true)
// and exit the program with the given status code.
func checkError(err error, exitStatus int, printUsage bool, cmd *cobra.Command) {
	if err == nil {
		return
	}

	fmt.Printf("Error: %v\n", err)
	if printUsage {
		cmd.Usage()
	}
	os.Exit(exitStatus)
}
