package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/cli"
	"github.com/lunjon/http/internal/format"
)

func main() {
	styler := format.NewStyler()
	cmd, err := cli.Build("v0.11.0")
	if err != nil {
		fmt.Println(styler.RedB("error:"), err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Println(styler.RedB("error:"), err)
		os.Exit(1)
	}
}
