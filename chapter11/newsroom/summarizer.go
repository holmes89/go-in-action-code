package newsroom

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// =============================================================================
// NEWSROOM MODULE - AI-powered article summarization
// =============================================================================

// Article represents the data needed for summarization.
// This is independent from the API's feed.Article to keep modules decoupled.
type Article struct {
	Title       string
	Description string
	Link        string
	FeedTitle   string
}

// =============================================================================
// CONFIGURATION
// =============================================================================

// Config holds configuration for the article summarizer.
type Config struct {
	OllamaURL string // Base URL for Ollama API
	Model     string // Model name to use
	MaxTokens int    // Maximum tokens in response
}

// DefaultConfig returns configuration with sensible defaults.
// Allows environment variable overrides for flexibility.
func DefaultConfig() Config {
	return Config{
		OllamaURL: getEnv("OLLAMA_URL", "http://localhost:11434"),
		Model:     getEnv("OLLAMA_MODEL", "llama2"),
		MaxTokens: 500,
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// =============================================================================
// ARTICLE SUMMARIZER - Real AI implementation
// =============================================================================

// ArticleSummarizer generates summaries using AI models via Ollama.
type ArticleSummarizer struct {
	llm    llms.Model
	config Config
}

// NewArticleSummarizer creates a new article summarizer with the given config.
// Returns an error if Ollama isn't available or model fails to load.
func NewArticleSummarizer(config Config) (*ArticleSummarizer, error) {
	// Create Ollama LLM client
	llm, err := ollama.New(
		ollama.WithModel(config.Model),
		ollama.WithServerURL(config.OllamaURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	return &ArticleSummarizer{
		llm:    llm,
		config: config,
	}, nil
}

// Summarize generates AI-powered summaries for articles.
// Handles both single-article summaries and multi-article news reports.
func (s *ArticleSummarizer) Summarize(ctx context.Context, articles []Article) (string, error) {
	if len(articles) == 0 {
		return "", fmt.Errorf("no articles to summarize")
	}

	// Build prompt based on article count
	var prompt string
	var maxTokens int

	if len(articles) == 1 {
		// Single article: concise summary
		article := articles[0]
		prompt = fmt.Sprintf(`Please provide a concise 2-3 sentence summary of this article:

Title: %s
Description: %s

Summary:`, article.Title, article.Description)
		maxTokens = 150
	} else {
		// Multiple articles: news report compilation
		var articleTexts []string
		for i, article := range articles {
			articleTexts = append(articleTexts, fmt.Sprintf(
				"Article %d:\nTitle: %s\nDescription: %s",
				i+1, article.Title, article.Description,
			))
		}

		prompt = fmt.Sprintf(`You are a news reporter. Compile these %d articles into a cohesive news report. 
Identify key themes and present them as a unified narrative.

%s

News Report:`, len(articles), strings.Join(articleTexts, "\n\n"))
		maxTokens = s.config.MaxTokens
	}

	// Call LLM with prompt
	response, err := llms.GenerateFromSinglePrompt(
		ctx,
		s.llm,
		prompt,
		llms.WithMaxTokens(maxTokens),
		llms.WithTemperature(0.7), // Balance creativity and consistency
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	return strings.TrimSpace(response), nil
}

// =============================================================================
// STUB SUMMARIZER - For development/testing without Ollama
// =============================================================================

// StubSummarizer provides a simple non-AI implementation.
// Useful for testing and development when Ollama isn't available.
type StubSummarizer struct{}

// NewStubSummarizer creates a new stub summarizer.
func NewStubSummarizer() *StubSummarizer {
	return &StubSummarizer{}
}

// Summarize generates a simple template-based summary.
func (s *StubSummarizer) Summarize(ctx context.Context, articles []Article) (string, error) {
	if len(articles) == 0 {
		return "", fmt.Errorf("no articles to summarize")
	}

	if len(articles) == 1 {
		return fmt.Sprintf(
			"[STUB] Summary of '%s': This article from %s discusses %s. For real AI summaries, install Ollama.",
			articles[0].Title,
			articles[0].FeedTitle,
			articles[0].Description,
		), nil
	}

	var titles []string
	for _, article := range articles {
		titles = append(titles, article.Title)
	}

	return fmt.Sprintf(
		"[STUB] News Report: Compiled %d articles covering: %s. For real AI summaries, install Ollama and configure OLLAMA_URL.",
		len(articles),
		strings.Join(titles, ", "),
	), nil
}

// This newsroom module demonstrates several key concepts:
//
// 1. Module Independence: Can evolve separately from the API
//    - Has its own go.mod
//    - Can be versioned independently
//    - Different release cadence
//
// 2. Adapter Pattern: Integrates external AI services
//    - Wraps Ollama/LangChain complexity
//    - Provides clean interface to consumers
//    - Can swap AI providers without changing consumers
//
// 3. Configuration Flexibility: Environment-based configuration
//    - Defaults for local development
//    - Override for different environments
//    - No hardcoded values
//
// 4. Graceful Degradation: Stub implementation when AI unavailable
//    - Development continues without Ollama
//    - Tests don't require external services
//    - Clear indication when using stub
//
// 5. Interface Satisfaction: Implicitly satisfies handlers.Summarizer
//    - No explicit declaration needed
//    - Duck typing through method signature
//    - Liskov Substitution Principle
//
// Future enhancements:
// - Support for other AI providers (OpenAI, Anthropic, Google)
// - Caching of summaries
// - Streaming responses for long-running summaries
// - Retry logic with exponential backoff
// - Circuit breaker for external service failures
// - Metrics and observability
