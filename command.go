package main

import (
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(s *state, cmd Command) error
}

func GetCommands() Commands {
	commands := Commands{
		Commands: map[string]func(s *state, cmd Command) error{},
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListAll)
	commands.register("agg", handlerFetch)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerGetFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", handlerFollowing)
	return commands
}

func (c *Commands) register(name string, f func(s *state, cmd Command) error) error {
	if _, ok := c.Commands[name]; ok {
		return fmt.Errorf("command already registered for: %s", name)
	}
	c.Commands[name] = f
	return nil
}

func (c *Commands) Run(s *state, cmd Command) error {
	if fn, ok := c.Commands[cmd.Name]; ok {
		if err := fn(s, cmd); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no command named: %s", cmd.Name)
	}
	return nil
}
