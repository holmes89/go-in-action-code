package store

import (
	"slices"
	"sync"

	feed "github.com/goinaction/go-news/api/internal"
)

// ArticleStore manages a thread-safe collection of articles with deduplication
// and sorting by publication date.
type ArticleStore struct {
	mu     sync.RWMutex
	cache  map[string]any
	sorted []*feed.Article
}

// NewArticleStore creates a new ArticleStore instance.
func NewArticleStore() *ArticleStore {
	return &ArticleStore{
		cache: make(map[string]any),
	}
}

// AddArticle adds a single article to the store, ensuring no duplicates.
// Articles are automatically sorted by publication date (newest first).
func (s *ArticleStore) AddArticle(article *feed.Article) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, ok := s.cache[article.Link]; ok {
		return
	}

	s.cache[article.Link] = nil
	s.sorted = append(s.sorted, article)
	slices.SortFunc(s.sorted, func(a, b *feed.Article) int {
		if a.Published == nil || b.Published == nil {
			return 0
		}
		return b.Published.Compare(*a.Published)
	})
}

// AddArticles adds multiple articles to the store in a batch operation.
// This is more efficient than calling AddArticle multiple times as it
// only sorts once after all articles are added.
func (s *ArticleStore) AddArticles(articles []*feed.Article) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	var dirty bool
	for _, article := range articles {
		if _, ok := s.cache[article.Link]; ok {
			continue
		}
		dirty = true
		s.cache[article.Link] = nil
		s.sorted = append(s.sorted, article)
	}
	
	if dirty {
		slices.SortFunc(s.sorted, func(a, b *feed.Article) int {
			if a.Published == nil || b.Published == nil {
				return 0
			}
			return b.Published.Compare(*a.Published)
		})
	}
}

// GetRecent returns the n most recent articles.
// If n is greater than the number of stored articles, all articles are returned.
func (s *ArticleStore) GetRecent(n int) []*feed.Article {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if n > len(s.sorted) {
		return s.sorted
	}
	return s.sorted[:n]
}
