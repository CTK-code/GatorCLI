package main

import (
	"context"
	"fmt"

	"github.com/CTK-code/GatorCLI/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd Command, user database.User) error) func(s *state, cmd Command) error {
	return func(s *state, cmd Command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error following - user may not exist?\n%s", err)
		}
		handler(s, cmd, user)
		return nil
	}
}
