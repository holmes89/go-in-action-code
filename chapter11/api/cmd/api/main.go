package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YOUR_USERNAME/go-news/api/internal/feed"
	"github.com/YOUR_USERNAME/go-news/api/internal/handlers"
	"github.com/YOUR_USERNAME/go-news/api/internal/reader"
	"github.com/YOUR_USERNAME/go-news/api/internal/store"
	"github.com/YOUR_USERNAME/go-news/newsroom"
)

// =============================================================================
// MAIN - Application entry point and dependency wiring
// =============================================================================

func main() {
	fmt.Println("Starting Go News API...")

	// Create core components with dependency injection
	// This demonstrates the composition root pattern - all wiring happens here

	// 1. Create storage - the single source of truth for articles
	articleStore := store.NewArticleStore()

	// 2. Create RSS reader with storage dependency
	rssReader := reader.NewRSSReader(articleStore)

	// 3. Create article handlers with read-only storage dependency
	articleHandlers := handlers.New(articleStore)

	// 4. Create AI summarizer with configuration
	config := newsroom.DefaultConfig()
	summarizer, err := newsroom.NewArticleSummarizer(config)
	if err != nil {
		// If Ollama isn't available, use stub implementation
		fmt.Printf("Warning: Failed to create AI summarizer: %v\n", err)
		fmt.Println("Using stub implementation. Install Ollama for real AI summaries.")
		summarizer = newsroom.NewStubSummarizer()
	}

	// 5. Create summary handlers with both dependencies
	summaryHandlers := handlers.NewSummaryHandlers(articleStore, summarizer)

	// Setup HTTP router
	mux := http.NewServeMux()

	// Let handlers register their own routes
	articleHandlers.RegisterRoutes(mux)
	summaryHandlers.RegisterRoutes(mux)

	// Add a root handler for documentation
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
  "service": "Go News API",
  "version": "1.0.0",
  "endpoints": {
    "GET /articles": "Fetch recent articles (supports ?count=N)",
    "GET /summary": "Generate AI news report (supports ?count=N)",
    "GET /": "This documentation"
  }
}`)
	})

	// Fetch initial feeds
	fmt.Println("Fetching initial feeds...")
	feeds := []string{
		"https://www.reddit.com/r/golang.rss",
		"https://go.dev/blog/feed.atom",
	}

	ctx := context.Background()
	for _, feedURL := range feeds {
		if _, err := rssReader.FetchFeed(ctx, feedURL); err != nil {
			fmt.Printf("Warning: Failed to fetch %s: %v\n", feedURL, err)
		}
	}

	// Start HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	// Start server in background
	go func() {
		fmt.Println("API server listening on http://localhost:8080")
		fmt.Println("\nTry these endpoints:")
		fmt.Println("  curl http://localhost:8080/")
		fmt.Println("  curl http://localhost:8080/articles?count=5")
		fmt.Println("  curl http://localhost:8080/summary?count=3")
		fmt.Println("\nPress Ctrl+C to stop")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Block until signal received
	<-done
	fmt.Println("\nShutting down gracefully...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}

// This main.go demonstrates several key patterns:
//
// 1. Composition Root: All dependency wiring happens here in main()
//    - No hidden dependencies or global state
//    - Easy to see what depends on what
//
// 2. Dependency Injection: Components receive dependencies through constructors
//    - Makes dependencies explicit
//    - Enables testing with mocks
//    - Prevents tight coupling
//
// 3. Interface-Based Design: Components depend on interfaces, not concrete types
//    - articleStore satisfies both feed.Storage and handlers.ArticleReader
//    - summarizer satisfies handlers.Summarizer
//    - Easy to swap implementations
//
// 4. Clean Architecture: Dependencies point inward
//    - main knows about everything
//    - domain (feed package) knows about nothing
//    - handlers depend on domain interfaces
//    - infrastructure (reader, store) implements domain interfaces
//
// 5. Graceful Shutdown: Server can finish in-flight requests
//    - Handles OS signals (Ctrl+C)
//    - Timeout prevents hanging
//    - Clean resource cleanup
//
// In production, you'd also add:
// - Configuration management (environment variables, config files)
// - Structured logging (zerolog, zap)
// - Metrics and observability (Prometheus, OpenTelemetry)
// - Health check endpoints
// - Background workers for periodic feed updates
// - Rate limiting and authentication
// - Database migrations
// - Feature flags
