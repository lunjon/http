package main

import (
	"fmt"

	"github.com/lunjon/http/internal/util"
)

func main() {
	b, err := util.OpenEditor("nvim")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
