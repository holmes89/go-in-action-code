# Chapter 11: Working with Larger Projects

A comprehensive demonstration of Go's organizational features: packages, modules, workspaces, clean architecture, domain-driven design, and multi-module development.

## Overview

This chapter builds a complete RSS feed reader with AI-powered summarization, demonstrating how to organize larger Go projects with multiple modules, clean architecture principles, and modern development workflows.

## Project Structure

```
chapter11/
├── go.work              # Workspace coordinating both modules
├── api/                 # RSS feed API module
│   ├── go.mod
│   ├── cmd/
│   │   └── api/
│   │       └── main.go  # HTTP server
│   └── internal/
│       ├── feed/        # Domain model (entities + interfaces)
│       ├── reader/      # RSS fetcher implementation
│       ├── store/       # In-memory storage
│       └── handlers/    # HTTP handlers
└── newsroom/            # AI summarization module
    ├── go.mod
    ├── summarizer.go    # AI-powered summarization
    └── cmd/
        └── newsroom/
            └── main.go  # Standalone CLI
```

## Key Concepts Demonstrated

### 1. **Modules and Packages**
- Multiple independent modules in a single repository
- Each module has its own `go.mod` and versioning
- Clean package organization with `internal/` for encapsulation
- `cmd/` directory for executables

### 2. **Domain-Driven Design (DDD)**
- Domain model (`feed` package) with entities and interfaces
- Domain defines "ports" (interfaces like `Fetcher` and `Storage`)
- Infrastructure provides "adapters" (implementations like `RSSReader` and `ArticleStore`)
- Clean separation between what and how

### 3. **Clean Architecture (Hexagonal Architecture)**
- Dependencies point inward toward the domain
- Domain has no dependencies on infrastructure
- Handlers depend on domain interfaces, not concrete implementations
- Easy to swap implementations without changing business logic

### 4. **Dependency Injection**
- Constructor injection makes dependencies explicit
- No hidden global state
- Easy to test with mock implementations
- Clear composition root in `main()`

### 5. **Interface-Based Design**
- Interfaces defined where they're consumed
- Implicit satisfaction (no `implements` keyword)
- Liskov Substitution Principle in action
- Flexible and testable code

### 6. **Workspaces**
- Local development of multiple modules
- Changes flow instantly between modules
- No version bumps during development
- Independent versioning for releases

### 7. **Testing Best Practices**
- Table-driven tests
- Mock implementations without frameworks
- Isolated unit tests
- Real HTTP testing with `httptest`

### 8. **AI Integration**
- Ollama for local LLM inference
- LangChain Go for AI frameworks
- Graceful degradation with stub implementation
- Environment-based configuration

## Architecture Layers

### Domain Layer (`internal/feed/`)
**Purpose**: Core business logic and contracts

```go
// Entities
type Article struct { ... }
type Feed struct { ... }

// Interfaces (Ports)
type Fetcher interface { ... }
type Storage interface { ... }
```

**Key Points**:
- No external dependencies
- Pure Go types
- Defines what the application needs
- Stable foundation

### Infrastructure Layer
**Purpose**: Concrete implementations of domain interfaces

**`internal/reader/`** - RSS fetching
- Implements `feed.Fetcher`
- Handles HTTP and XML parsing
- Converts external formats to domain types

**`internal/store/`** - Data storage
- Implements `feed.Storage`
- Thread-safe in-memory storage
- Could be swapped with PostgreSQL, MongoDB, etc.

### Presentation Layer (`internal/handlers/`)
**Purpose**: HTTP API and user interaction

- Defines its own narrow interfaces
- Converts between HTTP and domain types
- Thin layer - no business logic
- Testable with mocks

### Application Layer (`cmd/api/`)
**Purpose**: Composition root and startup

- Wires all dependencies together
- Creates concrete implementations
- Configures HTTP server
- Handles graceful shutdown

## Running the Project

### Prerequisites

```bash
# Install Go 1.22+
go version

# (Optional) Install Ollama for AI features
# Visit: https://ollama.ai
ollama pull llama2
```

### Run the API

```bash
cd api
go run ./cmd/api
```

The API starts on `http://localhost:8080` with endpoints:
- `GET /` - API documentation
- `GET /articles?count=N` - Fetch recent articles
- `GET /summary?count=N` - AI-generated news report

### Test the API

```bash
# Get API info
curl http://localhost:8080/

# Fetch 5 recent articles
curl http://localhost:8080/articles?count=5

# Generate news report from 3 articles
curl http://localhost:8080/summary?count=3
```

### Run Tests

```bash
# In api directory
cd api
go test ./...

# Run with coverage
go test -cover ./...

# Verbose output
go test -v ./internal/handlers
```

### Run Newsroom CLI

```bash
cd newsroom
go run ./cmd/newsroom
```

Demonstrates AI summarization standalone without the API.

## Configuration

Environment variables for customization:

```bash
# Ollama configuration
export OLLAMA_URL=http://localhost:11434
export OLLAMA_MODEL=llama2

# Run with custom config
go run ./cmd/api
```

If Ollama isn't available, the system automatically uses a stub implementation.

## Development Workflow

### 1. Initial Setup

```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/go-news
cd go-news

# Workspace is already configured with go.work
# Both modules are ready for development
```

### 2. Make Changes Across Modules

```bash
# Edit newsroom summarization logic
vim newsroom/summarizer.go

# Changes are immediately visible in API
# No need to publish or version bump

# Test API with new newsroom code
cd api
go run ./cmd/api
```

### 3. Test Modules Independently

```bash
# Test API module
cd api
go test ./...

# Test newsroom module
cd ../newsroom
go test ./...

# Test newsroom CLI
go run ./cmd/newsroom
```

### 4. Release Process

When ready to publish:

```bash
# Tag newsroom module
git tag -a newsroom/v0.2.0 -m "Add AI summarization"
git push origin newsroom/v0.2.0

# Update API's go.mod to use published version
cd api
go get github.com/YOUR_USERNAME/go-news/newsroom@v0.2.0

# Tag API module
git tag -a api/v1.5.2 -m "Add summary endpoint"
git push origin api/v1.5.2
```

## Design Patterns

### 1. Repository Pattern
`ArticleStore` encapsulates data access:
```go
type Storage interface {
    AddArticles([]*Article) error
    GetRecent(n int) []*Article
}
```

### 2. Adapter Pattern
`RSSReader` adapts RSS XML to domain types:
```go
// Convert external format → domain model
domainFeed := convertToFeed(rssFeed)
```

### 3. Dependency Inversion
High-level code depends on abstractions:
```go
// Handler depends on interface
type Handlers struct {
    articles ArticleReader // interface
}

// Store satisfies interface
store := NewArticleStore() // implements ArticleReader
```

### 4. Interface Segregation
Narrow interfaces for specific needs:
```go
// Handlers only need reading
type ArticleReader interface {
    GetRecent(n int) []*Article
}

// Not the full Storage interface
```

### 5. Composition Root
All wiring in one place:
```go
func main() {
    store := store.NewArticleStore()
    reader := reader.NewRSSReader(store)
    handlers := handlers.New(store)
    // ...
}
```

## Testing Strategy

### Unit Tests
Test components in isolation:
```go
// Mock dependency
mock := &mockArticleReader{...}

// Test with mock
handlers := New(mock)
```

### Table-Driven Tests
Multiple scenarios efficiently:
```go
tests := []struct {
    name string
    input string
    expected int
}{
    {"default", "", 10},
    {"custom", "?count=5", 5},
}
```

### HTTP Testing
Test real HTTP behavior:
```go
req := httptest.NewRequest("GET", "/articles", nil)
rec := httptest.NewRecorder()
mux.ServeHTTP(rec, req)
```

## Common Patterns

### Error Wrapping
```go
if err != nil {
    return fmt.Errorf("failed to fetch feed: %w", err)
}
```

### Context Propagation
```go
func (r *RSSReader) FetchFeed(ctx context.Context, url string) (*Feed, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    // Respects cancellation and timeouts
}
```

### Thread-Safe Storage
```go
func (s *ArticleStore) GetRecent(n int) []*Article {
    s.mu.RLock()
    defer s.mu.RUnlock()
    // ...
}
```

### Constructor Injection
```go
func NewRSSReader(storage feed.Storage) *RSSReader {
    return &RSSReader{storage: storage}
}
```

## Extending the Project

### Add Database Storage

```go
// Create PostgreSQL implementation
type PostgresStore struct {
    db *sql.DB
}

func (p *PostgresStore) AddArticles(articles []*feed.Article) error {
    // SQL INSERT logic
}

func (p *PostgresStore) GetRecent(n int) []*feed.Article {
    // SQL SELECT logic
}

// Swap in main()
store := NewPostgresStore(db) // Still satisfies feed.Storage
```

### Add Authentication

```go
// Create auth middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check auth token
        next.ServeHTTP(w, r)
    })
}

// Wrap mux
mux := http.NewServeMux()
server := &http.Server{
    Handler: authMiddleware(mux),
}
```

### Add Background Feed Updates

```go
func startFeedRefresher(ctx context.Context, reader feed.Fetcher, feeds []string) {
    ticker := time.NewTicker(15 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            for _, url := range feeds {
                reader.FetchFeed(ctx, url)
            }
        case <-ctx.Done():
            return
        }
    }
}
```

## Key Takeaways

1. **Organize by domain, not by type** - Group related functionality together
2. **Define interfaces where consumed** - Not where implemented
3. **Keep domain pure** - No infrastructure dependencies
4. **Inject dependencies explicitly** - Constructor injection over globals
5. **Use workspaces for multi-module development** - Fast iteration
6. **Test at boundaries** - Mock external dependencies
7. **Separate concerns clearly** - Domain, infrastructure, presentation
8. **Version modules independently** - Different release cadences

## Resources

- [Go Modules Reference](https://go.dev/ref/mod)
- [Go Workspace Tutorial](https://go.dev/doc/tutorial/workspaces)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Ollama Documentation](https://ollama.ai/docs)
- [LangChain Go](https://github.com/tmc/langchaingo)

## Troubleshooting

**Problem**: API can't find newsroom module
```bash
# Ensure workspace is configured
go work use ./api ./newsroom
```

**Problem**: Ollama connection failed
```bash
# Check Ollama is running
ollama serve

# Verify model is installed
ollama list

# Pull model if needed
ollama pull llama2
```

**Problem**: Tests failing
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Run tests with verbose output
go test -v ./...
```

**Problem**: Import cycle detected
- Review package dependencies
- Move shared code to common package
- Use interfaces to break cycles

## Summary

This chapter demonstrates production-ready Go project organization with:
- Clean architecture and domain-driven design
- Multi-module workspaces for independent versioning
- Dependency injection and interface-based design
- Comprehensive testing strategies
- AI integration with graceful degradation
- Modern development workflows

The patterns shown here scale from small projects to large distributed systems, providing a solid foundation for maintainable Go applications.
