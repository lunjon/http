package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/command"
)

func main() {
	const version = "0.8.0"
	http := command.Build(version)
	err := http.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
