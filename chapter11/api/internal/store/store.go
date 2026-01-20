package store

import (
	"slices"
	"sync"

	"github.com/YOUR_USERNAME/go-news/api/internal/feed"
)

// =============================================================================
// ARTICLE STORE - In-memory storage implementing feed.Storage
// =============================================================================

// Compile-time verification that ArticleStore implements feed.Storage
var _ feed.Storage = (*ArticleStore)(nil)

// ArticleStore provides thread-safe in-memory storage for articles.
// This demonstrates the Repository pattern - encapsulating data access
// behind an interface so the domain doesn't know about storage details.
type ArticleStore struct {
	mu       sync.RWMutex
	articles []*feed.Article
}

// NewArticleStore creates a new empty article store.
func NewArticleStore() *ArticleStore {
	return &ArticleStore{
		articles: make([]*feed.Article, 0),
	}
}

// AddArticles stores new articles in memory.
// Uses a write lock to ensure thread-safety during concurrent access.
func (s *ArticleStore) AddArticles(articles []*feed.Article) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Append new articles to the store
	s.articles = append(s.articles, articles...)

	// Sort by publication date (newest first) for efficient retrieval
	slices.SortFunc(s.articles, func(a, b *feed.Article) int {
		// Handle nil published dates
		if a.Published == nil && b.Published == nil {
			return 0
		}
		if a.Published == nil {
			return 1 // nil dates go to end
		}
		if b.Published == nil {
			return -1
		}

		// Sort descending (newest first)
		if a.Published.After(*b.Published) {
			return -1
		}
		if a.Published.Before(*b.Published) {
			return 1
		}
		return 0
	})

	return nil
}

// GetRecent returns the n most recent articles.
// Uses a read lock to allow concurrent reads while preventing writes.
func (s *ArticleStore) GetRecent(n int) []*feed.Article {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Handle edge cases
	if n <= 0 {
		return []*feed.Article{}
	}
	if n > len(s.articles) {
		n = len(s.articles)
	}

	// Return a copy to prevent external modification
	result := make([]*feed.Article, n)
	copy(result, s.articles[:n])

	return result
}

// In a production system, you might add methods like:
// - GetByFeed(feedTitle string) []*feed.Article
// - GetByDateRange(start, end time.Time) []*feed.Article
// - Delete(link string) error
// - Clear() error
//
// You could also implement feed.Storage with:
// - PostgreSQL using pgx
// - MongoDB using mongo-go-driver
// - SQLite using database/sql
//
// The beauty of the interface is that handlers and readers don't care
// which implementation you use - they all satisfy feed.Storage.
