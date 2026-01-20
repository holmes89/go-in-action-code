package feed

import (
	"context"
	"time"
)

// =============================================================================
// DOMAIN MODEL - Core entities independent of infrastructure
// =============================================================================

// Article represents a single article from an RSS feed.
// This is our core domain entity using only standard library types.
type Article struct {
	Title       string
	Description string
	Link        string
	Published   *time.Time
	FeedTitle   string
}

// Feed represents an RSS/Atom feed with its articles.
// This aggregates articles and provides feed-level metadata.
type Feed struct {
	Title       string
	Description string
	Link        string
	Articles    []*Article
}

// =============================================================================
// DOMAIN INTERFACES - Ports defining required behaviors
// =============================================================================

// Fetcher defines the behavior for fetching RSS feeds.
// This is a "port" in hexagonal architecture - it describes what we need
// without specifying how it's implemented. Implementations are "adapters".
type Fetcher interface {
	FetchFeed(ctx context.Context, url string) (*Feed, error)
}

// Storage defines how articles are persisted.
// Like Fetcher, this is an abstraction that can be satisfied by
// in-memory storage, databases, or any other implementation.
type Storage interface {
	AddArticles(articles []*Article) error
	GetRecent(n int) []*Article
}

// These interfaces demonstrate the Dependency Inversion Principle:
// High-level domain logic depends on abstractions (interfaces),
// not on low-level implementation details (concrete types).
//
// Benefits:
// - Domain remains stable while implementations can change
// - Easy to test with mock implementations
// - Clear contracts between layers
// - Prevents circular dependencies
