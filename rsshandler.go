package main

import (
	//"fmt"
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"html"
)

type RSSFeed struct {
	Channel struct {
		Title		string		`xml:"title"`
		Link		string		`xml:"link"`
		Description	string		`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	}	`xml:"channel"`
}

type RSSItem struct {
	Title		string	`xml:"title"`
	Link		string	`xml:"link"`
	Description	string	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}

func decodeFeed(f *RSSFeed) {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)

	for i, _ := range f.Channel.Item {
		f.Channel.Item[i].Title = html.UnescapeString(f.Channel.Item[i].Title)
		f.Channel.Item[i].Description = html.UnescapeString(f.Channel.Item[i].Description)
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	newFeed := RSSFeed{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {return &newFeed, err}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {return &newFeed, err}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {return &newFeed, err}

	err = xml.Unmarshal(data, &newFeed); if err != nil {
		return &newFeed, err
	}

	decodeFeed(&newFeed)

	return &newFeed, nil
}