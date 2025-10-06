package main

import (
	"time"
	"fmt"
	"context"
	"database/sql"
	"github.com/genus555/gator/internal/database"
	"github.com/lib/pq"
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

	for _, feed := range f.Channel.Item {
		if feed.Title == "" {continue}

		t, err := time.Parse(time.RFC1123Z, feed.PubDate)
		if err != nil {return err}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title:			feed.Title,
			Url:			feed.Link,
			Description:	sql.NullString{String: feed.Description, Valid: true},
			PublishedAt:	t,
			FeedID:			nextFeed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			fmt.Println("Error:", err)
		}
	}

	fmt.Printf("Saved posts in feed: \"%s\"\n", nextFeed.Name)

	return nil
}