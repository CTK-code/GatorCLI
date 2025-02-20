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

func handlerFetch(s *state, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("expected one argument")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s\n\n", timeBetweenRequests)
	fmt.Println("==========================================================================")
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
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

func handlerAddFeed(s *state, cmd Command, user database.User) error {
	fmt.Println(cmd.Args)
	if len(cmd.Args) < 2 {
		fmt.Println("expected two arguments")
		return errors.New("expected two arguments")
	}

	ctx := context.Background()
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

	// Add to following
	followArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	follow, err := s.db.CreateFeedFollow(ctx, followArgs)
	if err != nil {
		return fmt.Errorf("error following:\n%s", err)
	}

	fmt.Printf("New feed added:\n%v", feed)
	fmt.Printf("User: %s has followed %s\n",
		follow.UserName,
		follow.FeedName)
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

func handlerFollow(s *state, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("no arguments provided")
	}
	ctx := context.Background()

	feed, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not find feed:\n %s", err)
	}
	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	follow, err := s.db.CreateFeedFollow(ctx, args)
	if err != nil {
		return fmt.Errorf("error following:\n%s", err)
	}
	fmt.Printf("User: %s has followed %s\n",
		follow.UserName,
		follow.FeedName)
	return nil
}

func handlerFollowing(s *state, _ Command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeedFollowsForUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error fetcing following:\n%s", err)
	}
	fmt.Printf("Feeds for user: %s\n", s.Config.CurrentUserName)
	for _, feed := range feeds {
		fmt.Printf("*  %s\n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	args := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.Unfollow(ctx, args)
	if err != nil {
		return err
	}
	fmt.Printf("User %s has unfollowed %s\n", user.Name, feed.Name)
	return nil
}
