package api

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"strings"
)

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	reader := strings.NewReader("")
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}
	return &feed, nil
}

func (f *RSSFeed) UnescapeRSSFeed() {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Link = html.UnescapeString(f.Channel.Link)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)
	for _, i := range f.Channel.Item {
		i.Title = html.UnescapeString(i.Title)
		i.Link = html.UnescapeString(i.Link)
		i.Description = html.UnescapeString(i.Description)
		i.PubDate = html.UnescapeString(i.PubDate)
	}
}

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
