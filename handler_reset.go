package main

import (
	"context"
	"log"
)

func handlerReset(s *state, _ Command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
