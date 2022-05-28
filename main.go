package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/cli"
	"github.com/lunjon/http/internal/style"
)

func main() {
	cmd, err := cli.Build("v0.12.1")
	if err != nil {
		fmt.Println(style.RedB("error:"), err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Println(style.RedB("error:"), err)
		os.Exit(1)
	}
}
