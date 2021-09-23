package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/command"
)

func main() {
	cmd, err := command.Build("0.10.0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
