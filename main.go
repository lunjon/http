package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/command"
	"github.com/lunjon/http/style"
)

func main() {
	red := style.NewBuilder().Fg(style.Red).Bold(true).Build()

	cmd, err := command.Build("0.11.0")
	if err != nil {
		fmt.Println(red("error:"), err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Println(red("error:"), err)
		os.Exit(1)
	}
}
