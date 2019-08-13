package main

import (
	"fmt"
	"os"

	"github.com/lunjon/httpreq/internal/command"
)

func main() {
	httpreq := command.Build()
	err := httpreq.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
