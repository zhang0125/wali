package main

import (
	"fmt"

	"github.com/zhang0125/wali/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
