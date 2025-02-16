package main

import (
	"fmt"

	"github.com/CTK-code/GatorCLI/internal/config"
)

func main() {
	confData, err := config.Read()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(confData)
	}
}
