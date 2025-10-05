package main

import (
	"time"
	"fmt"
	"context"
	"database/sql"
	"github.com/genus555/gator/internal/database"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {return err}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt:		sql.NullTime{Time: time.Now(), Valid: true},
		ID:					nextFeed.ID,
	})
	if err != nil {return err}

	f, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {return err}

	fmt.Printf("Feed: \"%s\"\n", nextFeed.Name)
	fmt.Println("Posts:")

	for _, feed := range f.Channel.Item {
		fmt.Printf("   -\"%s\"\n", feed.Title)
	}

	return nil
}