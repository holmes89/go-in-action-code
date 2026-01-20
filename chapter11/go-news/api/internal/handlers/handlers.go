package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	feed "github.com/goinaction/go-news/api/internal"
	"github.com/goinaction/go-news/api/internal/store"
)

// FeedFetcher defines the interface for fetching feeds.
type FeedFetcher interface {
	FetchFeedSync(ctx context.Context, url string) ([]*feed.Article, error)
}

// ArticleProvider defines the interface for accessing stored articles.
type ArticleProvider interface {
	Store() *store.ArticleStore
}

// Handler manages HTTP request handling for the news aggregator.
type Handler struct {
	app       ArticleProvider
	fetcher   FeedFetcher
	feedURLs  []string
}

// NewHandler creates a new Handler with the given dependencies.
func NewHandler(app ArticleProvider, fetcher FeedFetcher, feedURLs []string) *Handler {
	return &Handler{
		app:      app,
		fetcher:  fetcher,
		feedURLs: feedURLs,
	}
}

// ArticlesHandler returns the most recent articles from the store.
// It accepts an optional "count" query parameter to limit results.
func (h *Handler) ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	n := 10
	if sizeStr := r.URL.Query().Get("count"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			n = size
		}
	}
	
	articles := h.app.Store().GetRecent(n)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// SyncHandler implements a fan-out/fan-in pattern to fetch all feeds
// concurrently and return aggregated results.
func (h *Handler) SyncHandler(w http.ResponseWriter, r *http.Request) {
	// Result structure for collecting feed fetch results
	type result struct {
		articles []*feed.Article
		err      error
	}
	
	resultsChan := make(chan result, len(h.feedURLs))
	var wg sync.WaitGroup
	
	// Fan-out: launch goroutines for each feed
	for _, feedURL := range h.feedURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			fetchedArticles, err := h.fetcher.FetchFeedSync(ctx, url)
			resultsChan <- result{articles: fetchedArticles, err: err}
		}(feedURL)
	}
	
	// Close results channel after all goroutines complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	// Fan-in: collect results from all goroutines
	var errors []string
	var articles []*feed.Article
	for res := range resultsChan {
		if res.err != nil {
			errors = append(errors, res.err.Error())
		} else {
			articles = append(articles, res.articles...)
		}
	}
	
	// Return errors if any occurred
	if len(errors) > 0 {
		w.WriteHeader(http.StatusPartialContent)
		response := map[string]interface{}{
			"articles": articles,
			"errors":   errors,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Return aggregated articles
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// HealthHandler provides a simple health check endpoint.
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "go-news-api",
	})
}

// StatusHandler returns information about the application status.
func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	articleCount := len(h.app.Store().GetRecent(1000000)) // Get all to count
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "running",
		"article_count": articleCount,
		"feed_count": len(h.feedURLs),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// LoggingMiddleware logs HTTP requests.
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		fmt.Printf(" - %v\n", time.Since(start))
	}
}
