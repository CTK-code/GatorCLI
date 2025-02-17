package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/CTK-code/GatorCLI/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func handlerFetch(_ *state, cmd Command) error {
	ctx := context.Background()
	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("Feed returned by fetch:\n%v\n", feed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected server response: %v", res.StatusCode)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var feed RSSFeed
	xml.Unmarshal(data, &feed)
	feed.unescape()
	return &feed, nil
}

func (feed *RSSFeed) unescape() {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {

		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}

func handlerAddFeed(s *state, cmd Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("expected two arguments")
	}
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error adding feed - user may not exist?\n%s", err)
	}
	args := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(ctx, args)
	if err != nil {
		return fmt.Errorf("error adding feed\n%s", err)
	}
	fmt.Printf("New feed added:\n%s", feed)
	return nil
}

func handlerGetFeeds(s *state, _ Command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("could not get feeds\n%s", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed Name: %s\n", feed.FeedName)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Added By: %s\n", feed.Username)
	}
	return nil
}
