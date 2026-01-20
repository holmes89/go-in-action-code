package internal

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
)

// ============================================================================
// Domain Models
// ============================================================================

// Article represents a news article from an RSS feed.
type Article struct {
	Title           string
	Description     string
	Content         string
	Link            string
	Published       *time.Time
	Updated         *time.Time
	Author          *Author
	GUID            string
	Categories      []string
	FeedTitle       string
	FeedDescription string
}

// Author represents the author of an article.
type Author struct {
	Name  string
	Email string
}

// Feed represents an RSS/Atom feed.
type Feed struct {
	Title       string
	Description string
	Link        string
	FeedLink    string
	Updated     *time.Time
	Published   *time.Time
	Language    string
	Articles    []*Article
}


// ============================================================================
// Interfaces
// ============================================================================

// Fetcher defines the interface for fetching and parsing RSS/Atom feeds.
type Fetcher interface {
	// FetchFeed retrieves and parses a feed from the given URL.
	FetchFeed(ctx context.Context, url string) (*Feed, error)
}

// ArticleSubmitter defines the interface for submitting articles to storage.
type ArticleSubmitter interface {
	// SubmitArticle submits a single article for storage.
	SubmitArticle(article *Article)
	
	// SubmitArticles submits multiple articles for storage.
	SubmitArticles(articles []*Article)
}

// ============================================================================
// GoFeed Adapter Implementation
// ============================================================================

// GoFeedAdapter adapts the gofeed library to our domain interface.
type GoFeedAdapter struct {
	parser *gofeed.Parser
}

// NewGoFeedAdapter creates a new gofeed adapter.
func NewGoFeedAdapter() *GoFeedAdapter {
	return &GoFeedAdapter{
		parser: gofeed.NewParser(),
	}
}

// FetchFeed retrieves and parses a feed using gofeed library.
func (a *GoFeedAdapter) FetchFeed(ctx context.Context, url string) (*Feed, error) {
	gfFeed, err := a.parser.ParseURLWithContext(ctx, url)
	if err != nil {
		return nil, err
	}
	
	return ConvertFeed(gfFeed), nil
}

// ConvertFeed converts a gofeed.Feed to our domain Feed model.
func ConvertFeed(gf *gofeed.Feed) *Feed {
	feed := &Feed{
		Title:       gf.Title,
		Description: gf.Description,
		Link:        gf.Link,
		FeedLink:    gf.FeedLink,
		Updated:     gf.UpdatedParsed,
		Published:   gf.PublishedParsed,
		Language:    gf.Language,
		Articles:    make([]*Article, 0, len(gf.Items)),
	}
	
	for _, item := range gf.Items {
		feed.Articles = append(feed.Articles, ConvertItem(item, feed))
	}
	
	return feed
}

// ConvertItem converts a gofeed.Item to our domain Article model.
func ConvertItem(item *gofeed.Item, feedData *Feed) *Article {
	article := &Article{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		Link:        item.Link,
		Published:   item.PublishedParsed,
		Updated:     item.UpdatedParsed,
		GUID:        item.GUID,
		Categories:  item.Categories,
	}
	
	if feedData != nil {
		article.FeedTitle = feedData.Title
		article.FeedDescription = feedData.Description
	}
	
	// Convert author if present
	if item.Author != nil {
		article.Author = &Author{
			Name:  item.Author.Name,
			Email: item.Author.Email,
		}
	}
	
	return article
}
