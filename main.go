package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CTK-code/GatorCLI/internal/config"
)

func main() {
	confData, err := config.Read()
	if err != nil {
		log.Fatal(err)
		return
	}
	programState := &state{
		Config: &confData,
	}
	commands := GetCommands()
	err = commands.Run(programState, argsToCommand())
	if err != nil {
		log.Fatal(err)
		return
	}

	confData, err = config.Read()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(confData)
}

func argsToCommand() Command {
	if len(os.Args) < 2 {
		return Command{}
	}
	// Ignore first arg which is the call to the program
	args := os.Args[1:]
	return Command{
		Name: args[0],
		Args: args[1:],
	}
}
