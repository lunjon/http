package main

import (
	"fmt"
	"os"

	"github.com/lunjon/httpreq/internal/command"
	"github.com/spf13/cobra"
)

const version = "v0.5.0"

func main() {
	v := &cobra.Command{
		Use:   `version`,
		Short: "Print the version of httpreq installed.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	httpreq := command.Build()
	httpreq.AddCommand(v)

	err := httpreq.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
