package main

import (
	"context"
	"fmt"
)

func handlerListAll(s *state, _ Command) error {
	ctx := context.Background()
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
