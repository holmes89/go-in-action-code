package reader

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/YOUR_USERNAME/go-news/api/internal/feed"
)

// =============================================================================
// RSS READER - Infrastructure adapter implementing feed.Fetcher
// =============================================================================

// Compile-time verification that RSSReader implements feed.Fetcher
var _ feed.Fetcher = (*RSSReader)(nil)

// RSSReader fetches and parses RSS feeds using a simplified RSS parser.
// In production, you'd typically use github.com/mmcdole/gofeed, but this
// demonstrates the adapter pattern: converting external formats to domain types.
type RSSReader struct {
	client  *http.Client
	storage feed.Storage // Dependency injection of storage interface
}

// NewRSSReader creates a new RSS reader with the given storage dependency.
// This is constructor injection - dependencies are explicit and testable.
func NewRSSReader(storage feed.Storage) *RSSReader {
	return &RSSReader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		storage: storage,
	}
}

// FetchFeed implements feed.Fetcher by fetching and parsing an RSS feed.
// This method demonstrates the full flow: fetch → parse → convert → store.
func (r *RSSReader) FetchFeed(ctx context.Context, url string) (*feed.Feed, error) {
	fmt.Printf("Fetching feed: %s\n", url)

	// Create HTTP request with context for cancellation support
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent to identify our application
	req.Header.Set("User-Agent", "Go-News-RSS-Reader/1.0")

	// Fetch the feed
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse as RSS 2.0 (simplified - real implementation would detect format)
	rssFeed, err := parseRSS(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS: %w", err)
	}

	// Convert from RSS structs to domain Feed type (adapter pattern)
	domainFeed := &feed.Feed{
		Title:       rssFeed.Channel.Title,
		Description: rssFeed.Channel.Description,
		Link:        rssFeed.Channel.Link,
		Articles:    make([]*feed.Article, 0, len(rssFeed.Channel.Items)),
	}

	// Convert each RSS item to a domain Article
	for _, item := range rssFeed.Channel.Items {
		article := &feed.Article{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			FeedTitle:   rssFeed.Channel.Title,
		}

		// Parse publication date if present
		if item.PubDate != "" {
			if pubTime, err := parseRFC822(item.PubDate); err == nil {
				article.Published = &pubTime
			}
		}

		domainFeed.Articles = append(domainFeed.Articles, article)
	}

	// Store articles using the injected storage dependency
	if err := r.storage.AddArticles(domainFeed.Articles); err != nil {
		return nil, fmt.Errorf("failed to store articles: %w", err)
	}

	fmt.Printf("Fetched %d articles from: %s\n", len(domainFeed.Articles), url)
	return domainFeed, nil
}

// =============================================================================
// RSS PARSING - XML structures for unmarshaling
// =============================================================================

// RSS structs represent the XML structure of an RSS 2.0 feed
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []item `xml:"item"`
}

type item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

// parseRSS unmarshals RSS XML into structured data
func parseRSS(data []byte) (*rss, error) {
	var rssFeed rss
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}
	return &rssFeed, nil
}

// parseRFC822 attempts to parse common RSS date formats.
// RSS 2.0 uses RFC 822, but feeds often vary in their date formatting.
func parseRFC822(dateStr string) (time.Time, error) {
	// Common RSS date formats to try in order
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
	}

	var lastErr error
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date %q: %w", dateStr, lastErr)
}
