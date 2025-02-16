package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login expects one argument username")
	}
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, cmd.Args[0])
	if err != nil {
		return err
	}
	if error := s.Config.SetUser(user.Name); error != nil {
		return error
	}
	fmt.Printf("User has been set to '%s'\n", user.Name)
	return nil
}
