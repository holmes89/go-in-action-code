package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/YOUR_USERNAME/go-news/api/internal/feed"
	"github.com/YOUR_USERNAME/go-news/newsroom"
)

// =============================================================================
// SUMMARY HANDLERS - AI-powered article summarization
// =============================================================================

// Summarizer defines the interface for article summarization.
// This demonstrates defining interfaces where they're consumed.
type Summarizer interface {
	Summarize(ctx context.Context, articles []newsroom.Article) (string, error)
}

// SummaryHandlers manages summary-related HTTP handlers.
type SummaryHandlers struct {
	articles   ArticleReader // For fetching articles
	summarizer Summarizer    // For AI-powered summarization
}

// NewSummaryHandlers creates handlers with summarization support.
// Both dependencies are injected through the constructor.
func NewSummaryHandlers(articles ArticleReader, summarizer Summarizer) *SummaryHandlers {
	return &SummaryHandlers{
		articles:   articles,
		summarizer: summarizer,
	}
}

// RegisterRoutes mounts summary handler routes on the provided mux.
func (h *SummaryHandlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/summary", h.newsReportHandler)
}

// newsReportHandler compiles recent articles into an AI-generated news report.
// Supports ?count=N query parameter to control number of articles.
func (h *SummaryHandlers) newsReportHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse count parameter with default of 5
	n := 5
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
			n = count
		}
	}

	// Fetch articles from storage
	articles := h.articles.GetRecent(n)

	if len(articles) == 0 {
		http.Error(w, "No articles available", http.StatusNotFound)
		return
	}

	// Convert from feed.Article to newsroom.Article
	// This demonstrates adapter pattern - converting between types
	newsroomArticles := make([]newsroom.Article, len(articles))
	for i, article := range articles {
		newsroomArticles[i] = newsroom.Article{
			Title:       article.Title,
			Description: article.Description,
			Link:        article.Link,
			FeedTitle:   article.FeedTitle,
		}
	}

	// Generate AI-powered summary
	summary, err := h.summarizer.Summarize(r.Context(), newsroomArticles)
	if err != nil {
		fmt.Printf("Failed to generate summary: %v\n", err)
		http.Error(w, "Failed to generate summary", http.StatusInternalServerError)
		return
	}

	// Return summary as JSON
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"article_count": len(articles),
		"summary":       summary,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// The beauty of this design:
//
// 1. Handler doesn't know if summarizer uses Ollama, ChatGPT, or a simple template
// 2. Easy to test with a mock summarizer
// 3. Can swap AI implementations without changing handler code
// 4. Clear separation between HTTP concerns and AI logic
// 5. Type conversion happens at the boundary, keeping domains clean
//
// This is the Open/Closed Principle in action:
// - Open for extension (new Summarizer implementations)
// - Closed for modification (handler code doesn't change)
