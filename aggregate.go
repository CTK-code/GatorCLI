package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CTK-code/GatorCLI/internal/database"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	args := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	}

	err = s.db.MarkFeedFetched(ctx, args)
	if err != nil {
		return err
	}
	rss, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return nil
	}
	fmt.Println("==========================================================================")
	fmt.Printf("Fetching %s From: %s", rss.Channel.Title, rss.Channel.Link)
	fmt.Println(rss.Channel.Description)
	fmt.Println("==========================================================================")
	for _, item := range rss.Channel.Item {
		fmt.Printf("  Fetching %s From: %s", item.Title, item.Link)
		fmt.Printf("  Published on %s", item.PubDate)
		fmt.Println("  " + item.Description)
	}

	return nil
}
