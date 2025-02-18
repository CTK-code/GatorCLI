package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CTK-code/GatorCLI/internal/database"
)

func handlerBrowse(s *state, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) >= 1 {
		var err error
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
	}

	ctx := context.Background()
	args := database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	}
	posts, err := s.db.GetPostsForUser(ctx, args)
	if err != nil {
		return err
	}
	fmt.Println("Browsing", len(posts), "posts.")

	for _, post := range posts {
		printPost(post)
	}
	return nil
}

func printPost(post database.Post) {
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Println(post.Description)
	fmt.Printf("\n\n")
}
