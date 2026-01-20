# Go News API

A concurrent RSS news aggregator demonstrating Go concurrency patterns including goroutines, channels, worker pools, and the fan-out/fan-in pattern.

## Project Structure

This project follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) conventions:

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go              # Main application coordinator
│   ├── feed/
│   │   └── feed.go             # Domain models, interfaces, and gofeed adapter
│   ├── handlers/
│   │   └── handlers.go         # HTTP request handlers
│   ├── store/
│   │   └── article_store.go    # Thread-safe article storage
│   └── worker/
│       └── worker.go           # Worker pool for concurrent feed fetching
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Features

- **Concurrent Feed Fetching**: Uses worker pools to fetch multiple RSS feeds concurrently
- **Thread-Safe Storage**: Mutex-protected article store with automatic deduplication
- **Channel-Based Communication**: Demonstrates Go channels for inter-goroutine communication
- **Fan-Out/Fan-In Pattern**: Synchronous endpoint that fetches all feeds concurrently
- **Context Support**: Proper timeout handling with Go contexts

## Building and Running

```bash
# Build the application
go build -o go-news ./cmd/api

# Run the application
./go-news
```

Or run directly:

```bash
go run ./cmd/api
```

## API Endpoints

### GET /articles?count=N
Returns the N most recent articles from the cache (default: 10).

```bash
curl http://localhost:8080/articles?count=5
```

### GET /sync
Fetches all configured feeds concurrently and returns aggregated results (fan-out/fan-in pattern).

```bash
curl http://localhost:8080/sync
```

### GET /health
Health check endpoint.

```bash
curl http://localhost:8080/health
```

### GET /status
Returns application status including article count and feed count.

```baDomain Models** (`pkg/models`): Clean domain entities (Article, Feed) independent of external libraries
2. **Feed Interface** (`pkg/feed`): Interface definitions for feed fetching and article submission
3. **GoFeed Adapter** (`internal/adapter`): Adapts the gofeed library to our domain models
4. **ArticleStore** (`internal/store`): Thread-safe storage with RWMutex for concurrent access
5. **Worker Pool** (`internal/worker`): Manages concurrent feed processing with configurable worker count
6. **App** (`internal/app`): Coordinates all components and manages channels
7
## Architecture

### Components

1. **ArticleStore** (`internal/store`): Thread-safe storage with RWMutex for concurrent access
2. **Worker Pool** (`internal/worker`): Manages concurrent feed fetching with configurable worker count
3. **App** (`internal/app`): Coordinates all components and manages channels
4. **Handlers** (`internal/handlers`): HTTP request handlers with middleware support

### Concurrency Patterns

- **Worker Pool**: Fixed number of goroutines processing feed URLs from a cha

### Architecture Patterns

- **Domain-Driven Design**: Clean domain models (Article, Feed) in `internal/feed`
- **Adapter Pattern**: GoFeedAdapter decouples external library (gofeed) from domain
- **Dependency Injection**: Interfaces allow for testability and flexibility
- **Interface Segregation**: Small, focused interfaces (`Fetcher`, `FeedProcessor`)nnel
- **Fan-Out/Fan-In**: Concurrent feed fetching with result aggregation
- **Channel Communication**: Non-blocking select for event handling
- **Mutex Synchronization**: Read-write locks for safe concurrent data access

## Dependencies

- [gofeed](https://github.com/mmcdole/gofeed): RSS/Atom feed parser
- [golang.org/x/sync](https://pkg.go.dev/golang.org/x/sync): Extended sync primitives

## Configuration

Edit the `feedURLs` variable in `cmd/api/main.go` to add or modify RSS feed sources:

```go
var feedURLs = []string{
    "https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
    "https://feeds.bbci.co.uk/news/rss.xml",
    // Add more feeds here
}
```

Worker pool size can be adjusted in `internal/app/app.go` (default: 3 workers).
