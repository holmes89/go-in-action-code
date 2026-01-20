# Chapter 9: Concurrency in Go

A comprehensive reference demonstrating Go's concurrency model through building an RSS reader application. This covers goroutines, channels, synchronization primitives, and common concurrency patterns.

## Quick Start

```bash
# Run all demonstrations
go run main.go

# Run with race detector (important!)
go run -race main.go

# To enable HTTP server, uncomment the section in main()
```

## What's Covered

### 1. **Concurrency Fundamentals**
- Concurrency vs Parallelism
- Goroutines with `go` keyword
- When to use (and not use) concurrency

### 2. **Synchronization Primitives**
- `sync.WaitGroup` - Wait for goroutines to complete
- `sync.Mutex` - Protect shared data (exclusive access)
- `sync.RWMutex` - Allow concurrent reads
- `errgroup` - Coordinate goroutines with error handling

### 3. **Context Package**
- Cancellation signals
- Timeouts and deadlines
- Request-scoped values
- Proper context propagation

### 4. **Channels**
- Unbuffered channels (synchronous)
- Buffered channels (asynchronous)
- Directional channels (read-only, write-only)
- Closing channels properly
- Range over channels

### 5. **select Statement**
- Multiplexing multiple channels
- Timeout patterns
- Non-blocking operations
- Default case usage

### 6. **Concurrency Patterns**
- Worker Pool (control max concurrency)
- Fan-Out/Fan-In (parallel work aggregation)
- Pipeline pattern (staged processing)

### 7. **Real-World Example**
- RSS Reader application
- Thread-safe article storage
- HTTP API with concurrent fetching
- Error handling in concurrent code

## Project Structure

```
chapter09/
├── main.go       # Complete implementation with demonstrations
├── go.mod        # Module definition
└── README.md     # This file
```

## Key Concepts Demonstrated

### Goroutines

```go
// Launch concurrent function
go fetchFeed("https://example.com/rss")

// With anonymous function
go func() {
    // Work here
}()
```

### WaitGroup Pattern

```go
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()
    // Work here
}()

wg.Wait() // Block until counter reaches 0
```

### Mutex for Shared Data

```go
type Store struct {
    mu   sync.RWMutex
    data map[string]string
}

func (s *Store) Set(key, val string) {
    s.mu.Lock()         // Exclusive access
    defer s.mu.Unlock()
    s.data[key] = val
}

func (s *Store) Get(key string) string {
    s.mu.RLock()        // Concurrent reads OK
    defer s.mu.RUnlock()
    return s.data[key]
}
```

### Context for Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := fetchWithContext(ctx, url)
```

### Basic Channels

```go
// Unbuffered (synchronous)
ch := make(chan string)
ch <- "hello"        // Send (blocks until received)
msg := <-ch          // Receive (blocks until sent)

// Buffered (asynchronous)
ch := make(chan int, 5)
ch <- 1              // Doesn't block until full
ch <- 2
ch <- 3
```

### Directional Channels

```go
// Send-only
func produce(out chan<- int) {
    out <- 42
}

// Receive-only
func consume(in <-chan int) {
    val := <-in
}
```

### select Statement

```go
select {
case msg := <-ch1:
    fmt.Println("Got:", msg)
case ch2 <- value:
    fmt.Println("Sent value")
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
default:
    fmt.Println("No activity (non-blocking)")
}
```

### Worker Pool Pattern

```go
type WorkerPool struct {
    jobs    chan Job
    results chan Result
}

func (p *WorkerPool) worker() {
    for job := range p.jobs {
        result := process(job)
        p.results <- result
    }
}

// Start N workers
for i := 0; i < workerCount; i++ {
    go p.worker()
}
```

### Fan-Out/Fan-In Pattern

```go
func fanOutFanIn(tasks []Task) []Result {
    results := make(chan Result, len(tasks))
    var wg sync.WaitGroup
    
    // Fan-out: distribute work
    for _, task := range tasks {
        wg.Add(1)
        go func(t Task) {
            defer wg.Done()
            result := process(t)
            results <- result
        }(task)
    }
    
    // Close channel when done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Fan-in: collect results
    var collected []Result
    for r := range results {
        collected = append(collected, r)
    }
    
    return collected
}
```

## RSS Reader Example

The main example builds an RSS reader that demonstrates:

### ArticleStore (Thread-Safe Storage)
```go
store := NewArticleStore()

// Safe from multiple goroutines
store.AddArticle(article)
store.AddArticles(articles)

// Concurrent reads OK
recent := store.GetRecent(10)
count := store.Count()
```

### Worker Pool for Feed Processing
```go
pool := NewWorkerPool(5, store) // 5 workers

// Submit jobs
pool.Submit("https://feed1.com/rss")
pool.Submit("https://feed2.com/rss")

// Workers process concurrently
```

### HTTP API Integration
```go
// GET /articles?count=10
// Returns recent articles

// GET /sync
// Fetches all feeds concurrently using fan-out/fan-in
```

## Best Practices

### ✓ DO

- **Keep it simple** - Avoid unnecessary complexity
- **Prefer channels** over shared memory for communication
- **Use Mutex** when sharing data is necessary
- **Pass context** as first parameter for cancellation
- **Limit goroutines** with worker pools
- **Close channels** properly (by sender)
- **Handle errors** explicitly in concurrent code
- **Test with** `go test -race` to catch race conditions
- **Track goroutines** with WaitGroup/errgroup

### ✗ DON'T

- **Launch unbounded** goroutines
- **Share data** without synchronization
- **Close from receiver** (only sender should close)
- **Use time.Sleep** for synchronization
- **Store contexts** in structs
- **Ignore errors** from goroutines
- **Use naked** `go` without tracking

## When to Use Concurrency

✅ **Good Use Cases:**
- Handling multiple network requests (web servers)
- I/O-bound operations (files, databases, network)
- Data pipelines and streaming
- Background tasks (email, notifications)
- Worker pools for job queues
- Real-time data feeds
- CPU-bound parallel computations

❌ **Poor Use Cases:**
- Simple, single-task programs
- Tasks requiring strict ordering
- When you're learning basics
- Overhead exceeds benefits
- Small, short-lived operations

## Common Pitfalls

### 1. Forgetting WaitGroup

```go
// ❌ Bad - goroutines may not finish
for _, url := range urls {
    go fetch(url)
}
// main exits immediately!

// ✅ Good
var wg sync.WaitGroup
for _, url := range urls {
    wg.Add(1)
    go func(u string) {
        defer wg.Done()
        fetch(u)
    }(url)
}
wg.Wait()
```

### 2. Race Conditions

```go
// ❌ Bad - race condition!
var counter int
for i := 0; i < 100; i++ {
    go func() {
        counter++ // UNSAFE!
    }()
}

// ✅ Good
var mu sync.Mutex
var counter int
for i := 0; i < 100; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

### 3. Channel Deadlock

```go
// ❌ Bad - deadlock!
ch := make(chan int)
ch <- 42 // Blocks forever (no receiver)

// ✅ Good - buffered or goroutine
ch := make(chan int, 1)
ch <- 42 // OK, buffered

// Or use goroutine
ch := make(chan int)
go func() {
    ch <- 42
}()
val := <-ch
```

### 4. Closing Closed Channels

```go
// ❌ Bad - panic!
ch := make(chan int)
close(ch)
close(ch) // PANIC!

// ✅ Good - only sender closes, once
ch := make(chan int)
go func() {
    defer close(ch)
    for _, val := range data {
        ch <- val
    }
}()
```

### 5. Goroutine Leaks

```go
// ❌ Bad - goroutines leak if never read
func leak() <-chan int {
    ch := make(chan int)
    go func() {
        // Blocks forever if ch never read
        ch <- compute()
    }()
    return ch
}

// ✅ Good - use context for cancellation
func noLeak(ctx context.Context) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        select {
        case ch <- compute():
        case <-ctx.Done():
            return
        }
    }()
    return ch
}
```

## Testing Concurrent Code

### Race Detector
```bash
# Always test with race detector
go test -race

go run -race main.go

go build -race
```

### Example Test
```go
func TestConcurrentAccess(t *testing.T) {
    store := NewArticleStore()
    var wg sync.WaitGroup
    
    // Launch 100 concurrent writes
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            article := Article{
                Title: fmt.Sprintf("Article %d", id),
                Link:  fmt.Sprintf("link%d", id),
            }
            store.AddArticle(article)
        }(i)
    }
    
    wg.Wait()
    
    if store.Count() != 100 {
        t.Errorf("Expected 100 articles, got %d", store.Count())
    }
}
```

## Performance Tips

1. **Profile First** - Don't guess, measure with pprof
2. **Limit Goroutines** - Use worker pools, not unlimited goroutines
3. **Buffered Channels** - Reduce blocking with appropriate buffer sizes
4. **RWMutex** - Use read locks for concurrent reads
5. **Batch Operations** - Reduce lock contention with batching
6. **sync.Pool** - Reuse objects to reduce GC pressure
7. **Avoid Contention** - Minimize shared state

## Debugging Concurrent Code

### Print Goroutine Stacks
```go
import "runtime/pprof"

pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
```

### Count Goroutines
```go
import "runtime"

fmt.Println("Goroutines:", runtime.NumGoroutine())
```

### Detect Deadlocks
```bash
# Go's runtime will detect obvious deadlocks
fatal error: all goroutines are asleep - deadlock!
```

## Further Reading

- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Advanced Go Concurrency Patterns](https://go.dev/blog/io2013-talk-concurrency)
- [Context Package](https://go.dev/blog/context)
- Book: "Concurrency in Go" by Katherine Cox-Buday
- Book: "Learn Concurrent Programming with Go" (Manning)

## Summary

✅ **Go's Concurrency Strengths:**
- Simple and powerful model (goroutines + channels)
- Lightweight goroutines (start with ~2KB stack)
- First-class channel support
- Built-in race detector
- Excellent standard library support

✅ **Key Patterns:**
- Worker pools for controlled concurrency
- Fan-out/fan-in for parallel aggregation
- Pipeline for staged processing
- Context for cancellation

✅ **Remember:**
- Start simple, add concurrency when needed
- Synchronize with channels or mutexes
- Always track goroutine completion
- Test with -race flag
- Profile before optimizing

Go's concurrency model makes it one of the best languages for building scalable, concurrent systems!
