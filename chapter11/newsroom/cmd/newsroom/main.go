package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/YOUR_USERNAME/go-news/newsroom"
)

// =============================================================================
// NEWSROOM CLI - Standalone application demonstrating summarization
// =============================================================================

func main() {
	fmt.Println("Go News - Article Summarizer")
	fmt.Println("=============================\n")

	// Load configuration from environment or use defaults
	config := newsroom.DefaultConfig()
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Ollama URL: %s\n", config.OllamaURL)
	fmt.Printf("  Model: %s\n", config.Model)
	fmt.Printf("  Max Tokens: %d\n\n", config.MaxTokens)

	// Create summarizer (will use stub if Ollama not available)
	summarizer, err := newsroom.NewArticleSummarizer(config)
	if err != nil {
		fmt.Printf("Warning: Failed to create AI summarizer: %v\n", err)
		fmt.Println("Using stub implementation instead.\n")
		stubSummarizer := newsroom.NewStubSummarizer()
		demonstrateWithSummarizer(stubSummarizer)
		return
	}

	demonstrateWithSummarizer(summarizer)
}

// demonstrateWithSummarizer runs demonstrations with any Summarizer implementation.
// This shows how interface-based design enables flexibility.
func demonstrateWithSummarizer(summarizer interface {
	Summarize(ctx context.Context, articles []newsroom.Article) (string, error)
}) {
	ctx := context.Background()

	// Create sample articles
	articles := []newsroom.Article{
		{
			Title:       "Go 1.23 Released with New Features",
			Description: "The Go team has released version 1.23 with improved performance, new standard library additions, and enhanced tooling support for larger projects.",
			Link:        "https://go.dev/blog/go1.23",
			FeedTitle:   "Go Blog",
		},
		{
			Title:       "Building Scalable APIs with Go",
			Description: "Learn best practices for designing RESTful APIs in Go, including proper error handling, middleware patterns, and graceful shutdown strategies.",
			Link:        "https://example.com/go-apis",
			FeedTitle:   "Go Weekly",
		},
		{
			Title:       "Domain-Driven Design in Go",
			Description: "An in-depth guide to implementing clean architecture and domain-driven design principles in Go applications, with practical examples and patterns.",
			Link:        "https://example.com/ddd-go",
			FeedTitle:   "Go Best Practices",
		},
	}

	// Demonstration 1: Single article summary
	fmt.Println("=== Single Article Summary ===\n")
	singleSummary, err := summarizer.Summarize(ctx, articles[:1])
	if err != nil {
		log.Fatalf("Failed to summarize single article: %v", err)
	}
	fmt.Printf("Article: %s\n", articles[0].Title)
	fmt.Printf("\nSummary:\n%s\n\n", singleSummary)

	time.Sleep(1 * time.Second) // Brief pause between requests

	// Demonstration 2: Multi-article news report
	fmt.Println("=== Multi-Article News Report ===\n")
	multiSummary, err := summarizer.Summarize(ctx, articles)
	if err != nil {
		log.Fatalf("Failed to generate news report: %v", err)
	}
	fmt.Printf("Compiled %d articles:\n", len(articles))
	for i, article := range articles {
		fmt.Printf("  %d. %s\n", i+1, article.Title)
	}
	fmt.Printf("\nNews Report:\n%s\n\n", multiSummary)

	fmt.Println("=============================")
	fmt.Println("Demonstration complete!")
	fmt.Println("\nTo use with real AI:")
	fmt.Println("  1. Install Ollama: https://ollama.ai")
	fmt.Println("  2. Pull a model: ollama pull llama2")
	fmt.Println("  3. Run this program again")
}

// This CLI demonstrates:
//
// 1. Module Independence: Can run standalone without the API
//    - newsroom is a separate module
//    - Can be developed and tested independently
//    - Could be published as a library
//
// 2. Configuration Management: Environment-based configuration
//    - Uses environment variables if set
//    - Falls back to sensible defaults
//    - Easy to configure for different environments
//
// 3. Graceful Degradation: Works with or without Ollama
//    - Attempts to use real AI
//    - Falls back to stub if unavailable
//    - Clear messaging about mode
//
// 4. Interface-Based Design: Works with any Summarizer
//    - Doesn't care about implementation details
//    - Could use OpenAI, Anthropic, or any provider
//    - Demonstrates Liskov Substitution Principle
//
// 5. Clear Output: User-friendly demonstrations
//    - Shows configuration
//    - Demonstrates both use cases
//    - Provides guidance for setup
//
// Usage examples:
//
//   # Use default configuration
//   go run ./cmd/newsroom
//
//   # Use custom Ollama URL
//   OLLAMA_URL=http://remote-server:11434 go run ./cmd/newsroom
//
//   # Use different model
//   OLLAMA_MODEL=mistral go run ./cmd/newsroom
//
// This standalone CLI is useful for:
// - Testing the newsroom module in isolation
// - Demonstrating AI capabilities to stakeholders
// - Debugging summarization logic
// - Benchmarking different models
// - Integration testing with different AI providers
