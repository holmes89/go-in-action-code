package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/YOUR_USERNAME/go-news/api/internal/feed"
)

// =============================================================================
// HTTP HANDLERS - Presentation layer for article API
// =============================================================================

// ArticleReader defines the read-only interface handlers need.
// This demonstrates Interface Segregation: handlers only depend on
// what they actually use (GetRecent), not the full feed.Storage interface.
type ArticleReader interface {
	GetRecent(n int) []*feed.Article
}

// Handlers manages HTTP request handlers with their dependencies.
type Handlers struct {
	articles ArticleReader
}

// New creates handlers with the given article reader dependency.
// Constructor injection makes dependencies explicit and testable.
func New(articles ArticleReader) *Handlers {
	return &Handlers{
		articles: articles,
	}
}

// RegisterRoutes mounts all handler routes on the provided mux.
// This pattern keeps route definitions co-located with their handlers
// and makes testing easier - tests can create their own mux.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/articles", h.articlesHandler)
}

// articlesHandler returns recent articles as JSON.
// Supports ?count=N query parameter to control number of articles returned.
func (h *Handlers) articlesHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse count parameter with default of 10
	n := 10
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
			n = count
		}
	}

	// Fetch articles from storage
	articles := h.articles.GetRecent(n)

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// In a production API, you'd add more handlers:
// - GET /articles/{id} - single article by ID
// - GET /feeds - list all feeds
// - POST /feeds - add a new feed to track
// - DELETE /feeds/{id} - stop tracking a feed
// - GET /health - health check endpoint
//
// Each handler would follow the same pattern:
// 1. Validate input
// 2. Call domain/storage interfaces
// 3. Transform to HTTP response
//
// The handler layer should be thin - business logic belongs in the domain.
