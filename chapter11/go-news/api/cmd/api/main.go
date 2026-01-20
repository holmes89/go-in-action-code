package main

import (
	"fmt"
	"net/http"

	"github.com/goinaction/go-news/api/internal/app"
	"github.com/goinaction/go-news/api/internal/handlers"
)

// feedURLs contains the list of RSS feeds to monitor.
var feedURLs = []string{
	"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
	"https://feeds.bbci.co.uk/news/rss.xml",
}

func main() {
	// Initialize the application
	application := app.NewApp()
	defer application.Close()

	// Submit initial feeds for background processing
	for _, feed := range feedURLs {
		application.SubmitFeed(feed)
	}

	// Set up HTTP handlers
	handler := handlers.NewHandler(application, application, feedURLs)
	
	http.HandleFunc("/articles", handlers.LoggingMiddleware(handler.ArticlesHandler))
	http.HandleFunc("/sync", handlers.LoggingMiddleware(handler.SyncHandler))
	http.HandleFunc("/health", handlers.LoggingMiddleware(handler.HealthHandler))
	http.HandleFunc("/status", handlers.LoggingMiddleware(handler.StatusHandler))
	
	// Start the HTTP server
	fmt.Println("HTTP server listening on :8080")
	fmt.Println("Endpoints:")
	fmt.Println("  GET /articles?count=10 - Get recent articles from cache")
	fmt.Println("  GET /sync             - Fetch all feeds and return aggregated results")
	fmt.Println("  GET /health           - Health check")
	fmt.Println("  GET /status           - Application status")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error:", err)
	}
}
