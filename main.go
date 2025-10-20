package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/cli"
	"github.com/lunjon/http/internal/style"
)

func main() {
	cmd, err := cli.Build("v0.14.0")
	if err != nil {
		fmt.Println(style.RedB.Render("error:"), err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Println(style.RedB.Render("error:"), err)
		os.Exit(1)
	}
}
