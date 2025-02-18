package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CTK-code/GatorCLI/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		publishedAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			return err
		}
		args := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}

		p, err := s.db.CreatePost(ctx, args)
		if err != nil {
			if _, ok := err.(*pq.Error); !ok {
				return err
			}
			fmt.Println("Post already added")
		} else {
			fmt.Printf("Post Added: %s", p.Title)
		}
	}

	return nil
}
