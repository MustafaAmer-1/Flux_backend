package main

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Language    string `xml:"language,omitempty"`
		PubDate     string `xml:"pubDate,omitempty"`
		LastBuild   string `xml:"lastBuildDate,omitempty"`
		Items       []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate,omitempty"`
	GUID        string `xml:"guid,omitempty"`
}

func fetchFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	var rssFeed RSSFeed
	err = decoder.Decode(&rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}
