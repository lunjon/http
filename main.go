package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/command"
	"github.com/lunjon/http/format"
)

func main() {
	cmd, err := command.Build("0.11.0")
	if err != nil {
		fmt.Println(format.RedB("error:"), err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Println(format.RedB("error:"), err)
		os.Exit(1)
	}
}
