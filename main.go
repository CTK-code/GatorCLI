package main

import (
	"fmt"
	"os"

	"github.com/CTK-code/GatorCLI/internal/config"
)

func main() {
	state := config.State{}
	confData, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(confData)
	}
	state.Config = &confData
	commands := config.GetCommands()
	err = commands.Run(&state, argsToCommand())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	confData, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(confData)
	}
}

func argsToCommand() config.Command {
	if len(os.Args) < 2 {
		return config.Command{}
	}
	// Ignore first arg which is the call to the program
	args := os.Args[1:]
	return config.Command{
		Name: args[0],
		Args: args[1:],
	}
}
