package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/YOUR_USERNAME/go-news/api/internal/feed"
	"github.com/YOUR_USERNAME/go-news/api/internal/handlers"
)

// =============================================================================
// HANDLER TESTS - Unit testing with dependency injection
// =============================================================================

// mockArticleReader is a test double for ArticleReader.
// This is a simple mock without using a mocking framework.
type mockArticleReader struct {
	articles []*feed.Article
}

func (m *mockArticleReader) GetRecent(n int) []*feed.Article {
	if n > len(m.articles) {
		n = len(m.articles)
	}
	return m.articles[:n]
}

// TestArticlesHandler demonstrates table-driven testing with dependency injection.
// Each test case is isolated and repeatable.
func TestArticlesHandler(t *testing.T) {
	// Test cases covering different scenarios
	tests := []struct {
		name           string
		queryString    string
		method         string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "default count",
			queryString:    "",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedCount:  10,
		},
		{
			name:           "custom count",
			queryString:    "?count=5",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedCount:  5,
		},
		{
			name:           "large count returns all",
			queryString:    "?count=100",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedCount:  10, // Only 10 articles in mock
		},
		{
			name:           "invalid method",
			queryString:    "",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock data - 10 test articles
			mockArticles := make([]*feed.Article, 10)
			now := time.Now()
			for i := 0; i < 10; i++ {
				mockArticles[i] = &feed.Article{
					Title:       "Article " + string(rune('A'+i)),
					Description: "Description " + string(rune('A'+i)),
					Link:        "https://example.com/" + string(rune('A'+i)),
					Published:   &now,
					FeedTitle:   "Test Feed",
				}
			}

			// Create mock reader
			mock := &mockArticleReader{articles: mockArticles}

			// Create handlers with mock dependency
			h := handlers.New(mock)

			// Create test HTTP mux and register routes
			mux := http.NewServeMux()
			h.RegisterRoutes(mux)

			// Create test request
			req := httptest.NewRequest(tt.method, "/articles"+tt.queryString, nil)
			rec := httptest.NewRecorder()

			// Serve the request
			mux.ServeHTTP(rec, req)

			// Verify status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// For successful requests, verify response body
			if tt.expectedStatus == http.StatusOK {
				var articles []*feed.Article
				if err := json.NewDecoder(rec.Body).Decode(&articles); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(articles) != tt.expectedCount {
					t.Errorf("expected %d articles, got %d", tt.expectedCount, len(articles))
				}
			}
		})
	}
}

// TestArticlesHandler_EmptyStore tests behavior with no articles.
func TestArticlesHandler_EmptyStore(t *testing.T) {
	// Create mock with no articles
	mock := &mockArticleReader{articles: []*feed.Article{}}

	h := handlers.New(mock)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var articles []*feed.Article
	if err := json.NewDecoder(rec.Body).Decode(&articles); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(articles) != 0 {
		t.Errorf("expected empty array, got %d articles", len(articles))
	}
}

// Benefits of this testing approach:
//
// 1. No external dependencies - tests run fast and reliably
// 2. Each test is isolated - no shared state between tests
// 3. Easy to add new test cases - just add to the table
// 4. Tests verify real HTTP behavior - using actual routes and mux
// 5. Mock is simple - no complex mocking framework needed
//
// In a production codebase, you'd also add:
// - Integration tests with real storage
// - Benchmark tests for performance
// - Fuzz tests for input validation
// - End-to-end tests with a real HTTP server
