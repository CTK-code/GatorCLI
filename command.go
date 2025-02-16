package main

import (
	"errors"
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

func handlerLogin(s *state, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login expects one argument username")
	}
	if err := s.Config.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	fmt.Printf("Username has been set to '%s'\n", cmd.Args[0])
	return nil
}
