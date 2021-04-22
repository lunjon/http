package main

import (
	"fmt"
	"os"

	"github.com/lunjon/httpreq/command"
)

func main() {
	const version = "0.8.1"
	httpreq := command.Build(version)
	err := httpreq.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
