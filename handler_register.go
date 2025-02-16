package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/CTK-code/GatorCLI/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("no arg passed for name")
	}
	ctx := context.Background()
	uuid := uuid.New()
	current_time := time.Now()

	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid,
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      cmd.Args[0],
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	s.Config.SetUser(user.Name)
	fmt.Printf("New user added: %v", user)

	return nil
}
