package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/CTK-code/GatorCLI/internal/config"
	"github.com/CTK-code/GatorCLI/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	confData, err := config.Read()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Set up the DB
	db, err := sql.Open("postgres", confData.DBURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()
	programState := &state{
		Config: &confData,
		db:     database.New(db),
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
	fmt.Printf("DBURL: %s\n", confData.DBURL)
	fmt.Printf("User: %s\n", confData.CurrentUserName)
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
