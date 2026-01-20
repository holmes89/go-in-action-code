package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"
)

// ============================================================================
// CHAPTER 9: CONCURRENCY IN GO
// ============================================================================
// This file demonstrates Go's concurrency model through building an RSS reader.
// We'll progressively show different concurrency concepts and patterns.
//
// Topics Covered:
// 1. Goroutines - Launching concurrent functions with `go`
// 2. sync.WaitGroup - Waiting for goroutines to complete
// 3. sync.Mutex/RWMutex - Protecting shared data
// 4. errgroup - Error handling in concurrent operations
// 5. context.Context - Cancellation and timeouts
// 6. Channels - Communication between goroutines
// 7. Buffered vs Unbuffered Channels
// 8. Directional Channels - Read-only and write-only
// 9. select Statement - Multiplexing channel operations
// 10. Worker Pool Pattern - Controlling concurrency
// 11. Fan-Out/Fan-In Pattern - Parallel work aggregation
//
// Key Concepts:
// - Concurrency vs Parallelism
// - Race conditions and how to avoid them
// - Channel communication patterns
// - Synchronization primitives
// - Best practices and when to use concurrency

// ============================================================================
// CONCURRENCY VS PARALLELISM
// ============================================================================
// CONCURRENCY: Managing multiple tasks at once (context switching)
// PARALLELISM: Actually executing multiple tasks simultaneously (multi-core)
//
// Example: Cooking multiple dishes
// - Concurrency: You switch between cooking chicken, doing laundry, sending email
// - Parallelism: Chicken bakes AND laundry runs at the exact same time
//
// Go makes concurrency easy with goroutines and channels.
// The Go runtime handles scheduling across available CPU cores.

// ============================================================================
// SIMULATED RSS FEED STRUCTURES
// ============================================================================
// For demonstration, we'll simulate RSS feeds instead of fetching real ones.
// In production, you'd use github.com/mmcdole/gofeed or similar.

// Article represents a news article from an RSS feed
type Article struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pub_date"`
	Source      string    `json:"source"`
}

// Feed represents an RSS feed source
type Feed struct {
	Title    string
	URL      string
	Articles []Article
}

// ============================================================================
// SIMULATED FEED FETCHER
// ============================================================================
// In a real application, this would make HTTP requests and parse XML/RSS.
// We simulate network delay and potential errors for demonstration.

// fetchFeedSimulated simulates fetching an RSS feed with network delay
func fetchFeedSimulated(url string) (Feed, error) {
	// Simulate network delay (100-500ms)
	time.Sleep(time.Duration(100+len(url)%400) * time.Millisecond)

	// Simulate occasional errors (10% failure rate)
	if len(url)%10 == 0 {
		return Feed{}, errors.New("simulated network error")
	}

	// Create simulated articles
	articles := []Article{
		{
			Title:       "Breaking News from " + url,
			Link:        url + "/article1",
			Description: "Important news story",
			PubDate:     time.Now().Add(-1 * time.Hour),
			Source:      url,
		},
		{
			Title:       "Tech Update from " + url,
			Link:        url + "/article2",
			Description: "Latest technology news",
			PubDate:     time.Now().Add(-2 * time.Hour),
			Source:      url,
		},
		{
			Title:       "World Events from " + url,
			Link:        url + "/article3",
			Description: "Global news coverage",
			PubDate:     time.Now().Add(-3 * time.Hour),
			Source:      url,
		},
	}

	return Feed{
		Title:    "Feed: " + url,
		URL:      url,
		Articles: articles,
	}, nil
}

// ============================================================================
// PART 1: BASIC GOROUTINES
// ============================================================================
// Goroutines are lightweight threads managed by the Go runtime.
// Launch with the `go` keyword: go functionName(args)

// printFeedTitle demonstrates a simple goroutine
func printFeedTitle(url string) {
	fmt.Printf("Fetching feed: %s\n", url)
	feed, err := fetchFeedSimulated(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	fmt.Printf("Fetched: %s - %s\n", url, feed.Title)
}

// demonstrateBasicGoroutines shows launching goroutines
func demonstrateBasicGoroutines() {
	fmt.Println("\n=== DEMONSTRATION 1: Basic Goroutines ===")

	fmt.Println("\n1. Synchronous (Sequential) Execution:")
	start := time.Now()
	printFeedTitle("https://news-source-1.com/rss")
	printFeedTitle("https://news-source-2.com/rss")
	fmt.Printf("Sequential time: %v\n", time.Since(start))

	fmt.Println("\n2. Asynchronous (Concurrent) Execution:")
	start = time.Now()
	// Launch goroutines - they run in background
	go printFeedTitle("https://news-source-3.com/rss")
	go printFeedTitle("https://news-source-4.com/rss")

	// Problem: main exits before goroutines finish!
	// We need synchronization
	time.Sleep(1 * time.Second) // Crude solution
	fmt.Printf("Concurrent time: %v (with sleep)\n", time.Since(start))

	fmt.Println("\nKey Points:")
	fmt.Println("• go keyword launches goroutines")
	fmt.Println("• Goroutines run concurrently with main")
	fmt.Println("• Need proper synchronization (not sleep!)")
}

// ============================================================================
// PART 2: SYNC.WAITGROUP - PROPER SYNCHRONIZATION
// ============================================================================
// WaitGroup waits for a collection of goroutines to finish.
// Three methods: Add(n), Done(), Wait()

// fetchFeedsWithWaitGroup demonstrates WaitGroup synchronization
func fetchFeedsWithWaitGroup(urls []string) {
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1) // Increment counter BEFORE launching goroutine

		go func(u string) {
			defer wg.Done() // Decrement counter when done

			fmt.Printf("Fetching: %s\n", u)
			feed, err := fetchFeedSimulated(u)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Success: %s\n", feed.Title)
		}(url) // Pass url as parameter to avoid closure issues
	}

	fmt.Println("Waiting for all fetches...")
	wg.Wait() // Block until counter reaches 0
	fmt.Println("All fetches complete!")
}

// demonstrateWaitGroup shows WaitGroup usage
func demonstrateWaitGroup() {
	fmt.Println("\n=== DEMONSTRATION 2: sync.WaitGroup ===")

	urls := []string{
		"https://nytimes.com/rss",
		"https://bbc.com/rss",
		"https://cnn.com/rss",
	}

	start := time.Now()
	fetchFeedsWithWaitGroup(urls)
	fmt.Printf("Total time: %v\n", time.Since(start))

	fmt.Println("\nWaitGroup Pattern:")
	fmt.Println("1. wg.Add(1) before launching goroutine")
	fmt.Println("2. defer wg.Done() inside goroutine")
	fmt.Println("3. wg.Wait() to block until all done")
}

// ============================================================================
// PART 3: ERRGROUP - ERROR HANDLING IN CONCURRENT CODE
// ============================================================================
// errgroup.Group provides synchronization AND error collection.
// If any goroutine returns an error, Wait() returns that error.

// Note: In production, use "golang.org/x/sync/errgroup"
// For this demo, we'll simulate its behavior

// SimpleErrGroup is a simplified version of errgroup.Group
type SimpleErrGroup struct {
	wg  sync.WaitGroup
	mu  sync.Mutex
	err error
}

func (g *SimpleErrGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.mu.Lock()
			if g.err == nil {
				g.err = err // Store first error
			}
			g.mu.Unlock()
		}
	}()
}

func (g *SimpleErrGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

// fetchFeedWithError wraps fetch and returns error
func fetchFeedWithError(url string) (Feed, error) {
	fmt.Printf("Fetching: %s\n", url)
	feed, err := fetchFeedSimulated(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return Feed{}, err
	}
	fmt.Printf("Success: %s\n", feed.Title)
	return feed, nil
}

// demonstrateErrGroup shows error handling with errgroup
func demonstrateErrGroup() {
	fmt.Println("\n=== DEMONSTRATION 3: Error Handling (errgroup) ===")

	urls := []string{
		"https://source1.com/rss",
		"https://source2.com/rss",
		"https://source3000000000.com/rss", // Will fail (length % 10 == 0)
	}

	var eg SimpleErrGroup
	var feeds []Feed
	var mu sync.Mutex // Protect feeds slice

	for _, url := range urls {
		u := url
		eg.Go(func() error {
			feed, err := fetchFeedWithError(u)
			if err != nil {
				return err
			}
			mu.Lock()
			feeds = append(feeds, feed)
			mu.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("\nAt least one fetch failed: %v\n", err)
	} else {
		fmt.Println("\nAll fetches succeeded!")
	}

	fmt.Printf("Successfully fetched %d feeds\n", len(feeds))
	fmt.Println("\nKey Points:")
	fmt.Println("• errgroup.Go() launches goroutines")
	fmt.Println("• Return error from each goroutine")
	fmt.Println("• Wait() returns first error encountered")
}

// ============================================================================
// PART 4: ARTICLE STORE - SHARED STATE WITH MUTEX
// ============================================================================
// When multiple goroutines access shared data, use Mutex for synchronization.
// Mutex = Mutually Exclusive lock
// RWMutex = Read/Write lock (allows concurrent reads)

// ArticleStore safely stores articles from multiple goroutines
type ArticleStore struct {
	mu     sync.RWMutex      // Protects cache and sorted
	cache  map[string]bool   // Deduplicate by link
	sorted []Article         // Articles sorted by date
}

// NewArticleStore creates an initialized store
func NewArticleStore() *ArticleStore {
	return &ArticleStore{
		cache:  make(map[string]bool),
		sorted: make([]Article, 0),
	}
}

// AddArticle safely adds a single article (thread-safe)
func (s *ArticleStore) AddArticle(article Article) {
	s.mu.Lock()         // Write lock - exclusive access
	defer s.mu.Unlock() // Always unlock

	// Check if already exists
	if s.cache[article.Link] {
		return
	}

	// Add to cache and sorted list
	s.cache[article.Link] = true
	s.sorted = append(s.sorted, article)

	// Sort by publication date (newest first)
	slices.SortFunc(s.sorted, func(a, b Article) int {
		if a.PubDate.After(b.PubDate) {
			return -1
		}
		if a.PubDate.Before(b.PubDate) {
			return 1
		}
		return 0
	})
}

// AddArticles safely adds multiple articles
func (s *ArticleStore) AddArticles(articles []Article) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hasNew := false
	for _, article := range articles {
		if !s.cache[article.Link] {
			s.cache[article.Link] = true
			s.sorted = append(s.sorted, article)
			hasNew = true
		}
	}

	// Only sort if we added new articles
	if hasNew {
		slices.SortFunc(s.sorted, func(a, b Article) int {
			if a.PubDate.After(b.PubDate) {
				return -1
			}
			if a.PubDate.Before(b.PubDate) {
				return 1
			}
			return 0
		})
	}
}

// GetRecent safely retrieves the n most recent articles
func (s *ArticleStore) GetRecent(n int) []Article {
	s.mu.RLock()         // Read lock - allows concurrent reads
	defer s.mu.RUnlock() // Always unlock

	if n > len(s.sorted) {
		n = len(s.sorted)
	}

	// Return a copy to prevent external modification
	result := make([]Article, n)
	copy(result, s.sorted[:n])
	return result
}

// Count returns the total number of unique articles
func (s *ArticleStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sorted)
}

// demonstrateMutex shows thread-safe shared state
func demonstrateMutex() {
	fmt.Println("\n=== DEMONSTRATION 4: Mutex (Thread-Safe Store) ===")

	store := NewArticleStore()
	var wg sync.WaitGroup

	urls := []string{
		"https://source1.com/rss",
		"https://source2.com/rss",
		"https://source3.com/rss",
	}

	// Launch goroutines to fetch and store articles
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			feed, err := fetchFeedSimulated(u)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Safe to call from multiple goroutines
			store.AddArticles(feed.Articles)
			fmt.Printf("Added articles from %s\n", u)
		}(url)
	}

	wg.Wait()

	fmt.Printf("\nTotal articles stored: %d\n", store.Count())
	fmt.Println("Recent articles:")
	for i, article := range store.GetRecent(5) {
		fmt.Printf("%d. %s (from %s)\n", i+1, article.Title, article.Source)
	}

	fmt.Println("\nMutex Patterns:")
	fmt.Println("• Lock() for exclusive access (writes)")
	fmt.Println("• RLock() for concurrent reads")
	fmt.Println("• Always defer Unlock()")
	fmt.Println("• Keep critical sections small")
}

// ============================================================================
// PART 5: CONTEXT - CANCELLATION AND TIMEOUTS
// ============================================================================
// Context carries cancellation signals, deadlines, and request-scoped values.
// Use context.WithTimeout, context.WithCancel, context.WithDeadline

// fetchFeedWithContext respects context cancellation
func fetchFeedWithContext(ctx context.Context, url string) (Feed, error) {
	// Simulate checking context before expensive operation
	select {
	case <-ctx.Done():
		return Feed{}, ctx.Err() // Cancelled or timed out
	default:
		// Continue
	}

	fmt.Printf("Fetching: %s\n", url)

	// Simulate long operation with context checking
	done := make(chan struct{})
	var feed Feed
	var err error

	go func() {
		feed, err = fetchFeedSimulated(url)
		close(done)
	}()

	select {
	case <-done:
		// Completed successfully
		if err != nil {
			return Feed{}, err
		}
		fmt.Printf("Success: %s\n", feed.Title)
		return feed, nil
	case <-ctx.Done():
		// Context cancelled or timed out
		return Feed{}, ctx.Err()
	}
}

// demonstrateContext shows context usage
func demonstrateContext() {
	fmt.Println("\n=== DEMONSTRATION 5: Context (Cancellation & Timeout) ===")

	// Example 1: Timeout context
	fmt.Println("\n1. With 2-second timeout:")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Always call cancel to release resources

	var wg sync.WaitGroup
	urls := []string{
		"https://fast-source.com/rss",
		"https://medium-source.com/rss",
	}

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			_, err := fetchFeedWithContext(ctx, u)
			if err != nil {
				fmt.Printf("Failed %s: %v\n", u, err)
			}
		}(url)
	}

	wg.Wait()

	// Example 2: Manual cancellation
	fmt.Println("\n2. With manual cancellation:")
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Cancelling all operations...")
		cancel2()
	}()

	_, err := fetchFeedWithContext(ctx2, "https://slow-source.com/rss")
	if err != nil {
		fmt.Printf("Fetch cancelled: %v\n", err)
	}

	fmt.Println("\nContext Best Practices:")
	fmt.Println("• Pass context as first parameter")
	fmt.Println("• Always defer cancel()")
	fmt.Println("• Check ctx.Done() in long operations")
	fmt.Println("• Don't store contexts in structs")
}

// ============================================================================
// PART 6: CHANNELS - COMMUNICATION BETWEEN GOROUTINES
// ============================================================================
// Channels are typed conduits for sending/receiving values between goroutines.
// Create with make(), send with <-, receive with <-

// demonstrateBasicChannels shows channel fundamentals
func demonstrateBasicChannels() {
	fmt.Println("\n=== DEMONSTRATION 6: Channels (Communication) ===")

	// Example 1: Unbuffered channel (synchronous)
	fmt.Println("\n1. Unbuffered Channel (size 0):")
	ch := make(chan string)

	go func() {
		fmt.Println("   Goroutine: Sending message...")
		ch <- "Hello from goroutine!" // Blocks until received
		fmt.Println("   Goroutine: Message sent")
	}()

	time.Sleep(100 * time.Millisecond) // Show goroutine blocks
	fmt.Println("   Main: Receiving message...")
	msg := <-ch // Blocks until sent
	fmt.Printf("   Main: Received: %s\n", msg)

	// Example 2: Buffered channel (asynchronous)
	fmt.Println("\n2. Buffered Channel (size 3):")
	buffered := make(chan int, 3)

	buffered <- 1 // Doesn't block
	buffered <- 2 // Doesn't block
	buffered <- 3 // Doesn't block
	fmt.Println("   Sent 3 values without blocking")

	// buffered <- 4 // Would block here (buffer full)

	fmt.Printf("   Received: %d\n", <-buffered)
	fmt.Printf("   Received: %d\n", <-buffered)
	fmt.Printf("   Received: %d\n", <-buffered)

	// Example 3: Channel for results
	fmt.Println("\n3. Channel for Concurrent Results:")
	results := make(chan Feed, 3)

	urls := []string{"https://a.com/rss", "https://b.com/rss", "https://c.com/rss"}
	for _, url := range urls {
		go func(u string) {
			feed, _ := fetchFeedSimulated(u)
			results <- feed // Send result
		}(url)
	}

	// Collect results
	for i := 0; i < len(urls); i++ {
		feed := <-results
		fmt.Printf("   Got feed: %s\n", feed.Title)
	}

	fmt.Println("\nChannel Concepts:")
	fmt.Println("• make(chan T) creates unbuffered channel")
	fmt.Println("• make(chan T, n) creates buffered channel")
	fmt.Println("• ch <- val sends value")
	fmt.Println("• val := <-ch receives value")
	fmt.Println("• Unbuffered = synchronous handoff")
	fmt.Println("• Buffered = asynchronous up to capacity")
}

// ============================================================================
// PART 7: DIRECTIONAL CHANNELS - READ-ONLY AND WRITE-ONLY
// ============================================================================
// Channels can be restricted to send-only or receive-only.
// This provides type safety and makes intent clear.

// producer sends articles to a write-only channel
func producer(articles chan<- Article, source string) {
	for i := 1; i <= 3; i++ {
		article := Article{
			Title:   fmt.Sprintf("Article %d from %s", i, source),
			Link:    fmt.Sprintf("%s/article%d", source, i),
			PubDate: time.Now(),
			Source:  source,
		}
		articles <- article // Can only send
	}
}

// consumer receives articles from a read-only channel
func consumer(articles <-chan Article, id int) {
	for article := range articles { // Can only receive
		fmt.Printf("   Consumer %d: %s\n", id, article.Title)
	}
}

// demonstrateDirectionalChannels shows channel direction
func demonstrateDirectionalChannels() {
	fmt.Println("\n=== DEMONSTRATION 7: Directional Channels ===")

	articles := make(chan Article, 10)
	var producerWg sync.WaitGroup

	// Launch producers (write-only channels)
	producerWg.Add(2)
	go func() {
		defer producerWg.Done()
		producer(articles, "Source-A")
	}()
	go func() {
		defer producerWg.Done()
		producer(articles, "Source-B")
	}()

	// Close channel after all producers finish
	go func() {
		producerWg.Wait()
		close(articles)
	}()

	// Launch consumers (read-only channels)
	var consumerWg sync.WaitGroup
	for i := 1; i <= 2; i++ {
		consumerWg.Add(1)
		go func(id int) {
			defer consumerWg.Done()
			consumer(articles, id)
		}(i)
	}

	consumerWg.Wait()

	fmt.Println("\nDirectional Channel Syntax:")
	fmt.Println("• chan<- T : send-only (write-only)")
	fmt.Println("• <-chan T : receive-only (read-only)")
	fmt.Println("• Provides type safety and intent")
}

// ============================================================================
// PART 8: SELECT STATEMENT - MULTIPLEXING CHANNELS
// ============================================================================
// select lets a goroutine wait on multiple channel operations.
// Like switch, but for channels. Proceeds with whichever is ready first.

// demonstrateSelect shows select statement patterns
func demonstrateSelect() {
	fmt.Println("\n=== DEMONSTRATION 8: select Statement ===")

	// Example 1: Timeout pattern
	fmt.Println("\n1. Timeout Pattern:")
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "data"
	}()

	select {
	case msg := <-ch:
		fmt.Printf("   Received: %s\n", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("   Timeout! No data received")
	}

	// Example 2: Multiple channels
	fmt.Println("\n2. Multiple Channels:")
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "from ch2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("   Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("   Received: %s\n", msg2)
		}
	}

	// Example 3: Non-blocking receive
	fmt.Println("\n3. Non-Blocking Receive:")
	ch3 := make(chan int)

	select {
	case val := <-ch3:
		fmt.Printf("   Received: %d\n", val)
	default:
		fmt.Println("   No data available (non-blocking)")
	}

	fmt.Println("\nselect Patterns:")
	fmt.Println("• Wait on multiple channels")
	fmt.Println("• Proceed with first ready channel")
	fmt.Println("• Use default for non-blocking")
	fmt.Println("• time.After for timeouts")
}

// ============================================================================
// PART 9: WORKER POOL PATTERN
// ============================================================================
// Worker pools control concurrency by using a fixed number of goroutines
// that process tasks from a shared queue (channel).

// Worker processes feed URLs from a channel
type Worker struct {
	id      int
	jobs    <-chan string  // Read-only channel of URLs
	results chan<- Feed    // Write-only channel of results
	store   *ArticleStore
}

// run processes jobs until channel is closed
func (w *Worker) run() {
	for url := range w.jobs {
		fmt.Printf("Worker %d: Processing %s\n", w.id, url)
		feed, err := fetchFeedSimulated(url)
		if err != nil {
			fmt.Printf("Worker %d: Error: %v\n", w.id, err)
			continue
		}
		w.store.AddArticles(feed.Articles)
		w.results <- feed
	}
	fmt.Printf("Worker %d: Finished\n", w.id)
}

// WorkerPool manages multiple workers
type WorkerPool struct {
	workers []*Worker
	jobs    chan string
	results chan Feed
	store   *ArticleStore
}

// NewWorkerPool creates a pool with n workers
func NewWorkerPool(size int, store *ArticleStore) *WorkerPool {
	jobs := make(chan string, 100)
	results := make(chan Feed, 100)

	pool := &WorkerPool{
		workers: make([]*Worker, size),
		jobs:    jobs,
		results: results,
		store:   store,
	}

	// Start workers
	for i := 0; i < size; i++ {
		worker := &Worker{
			id:      i + 1,
			jobs:    jobs,
			results: results,
			store:   store,
		}
		pool.workers[i] = worker
		go worker.run()
	}

	return pool
}

// Submit adds a job to the pool
func (p *WorkerPool) Submit(url string) {
	p.jobs <- url
}

// Close shuts down the pool
func (p *WorkerPool) Close() {
	close(p.jobs)
}

// demonstrateWorkerPool shows worker pool pattern
func demonstrateWorkerPool() {
	fmt.Println("\n=== DEMONSTRATION 9: Worker Pool Pattern ===")

	store := NewArticleStore()
	pool := NewWorkerPool(3, store) // 3 workers

	urls := []string{
		"https://feed1.com/rss",
		"https://feed2.com/rss",
		"https://feed3.com/rss",
		"https://feed4.com/rss",
		"https://feed5.com/rss",
		"https://feed6.com/rss",
	}

	// Submit jobs
	fmt.Println("Submitting jobs to pool...")
	for _, url := range urls {
		pool.Submit(url)
	}

	// Collect results
	go func() {
		for i := 0; i < len(urls); i++ {
			feed := <-pool.results
			fmt.Printf("Result: %s\n", feed.Title)
		}
	}()

	time.Sleep(1 * time.Second)
	pool.Close()

	fmt.Printf("\nProcessed %d articles total\n", store.Count())

	fmt.Println("\nWorker Pool Benefits:")
	fmt.Println("• Control max concurrency")
	fmt.Println("• Reuse goroutines (efficient)")
	fmt.Println("• Fair work distribution")
	fmt.Println("• Backpressure with buffered jobs channel")
}

// ============================================================================
// PART 10: FAN-OUT/FAN-IN PATTERN
// ============================================================================
// Fan-out: Distribute work to multiple goroutines
// Fan-in: Collect results from multiple goroutines

// fanOutFanIn demonstrates the pattern
func fanOutFanIn(urls []string) ([]Feed, error) {
	// Fan-out: Launch goroutine for each URL
	results := make(chan Feed, len(urls))
	errors := make(chan error, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			feed, err := fetchFeedSimulated(u)
			if err != nil {
				errors <- err
				return
			}
			results <- feed
		}(url)
	}

	// Close channels after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Fan-in: Collect all results
	var feeds []Feed
	var errs []error

	// Collect from both channels until closed
	for results != nil || errors != nil {
		select {
		case feed, ok := <-results:
			if !ok {
				results = nil // Channel closed
			} else {
				feeds = append(feeds, feed)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil // Channel closed
			} else {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return feeds, fmt.Errorf("encountered %d errors", len(errs))
	}

	return feeds, nil
}

// demonstrateFanOutFanIn shows the pattern
func demonstrateFanOutFanIn() {
	fmt.Println("\n=== DEMONSTRATION 10: Fan-Out/Fan-In Pattern ===")

	urls := []string{
		"https://source1.com/rss",
		"https://source2.com/rss",
		"https://source3.com/rss",
		"https://source4.com/rss",
	}

	fmt.Println("Fan-out: Launching parallel fetches...")
	start := time.Now()

	feeds, err := fanOutFanIn(urls)

	fmt.Printf("Fan-in: Collected %d feeds in %v\n", len(feeds), time.Since(start))

	if err != nil {
		fmt.Printf("Errors occurred: %v\n", err)
	}

	for _, feed := range feeds {
		fmt.Printf("  - %s (%d articles)\n", feed.Title, len(feed.Articles))
	}

	fmt.Println("\nFan-Out/Fan-In Pattern:")
	fmt.Println("• Fan-out: Parallel work distribution")
	fmt.Println("• Fan-in: Result aggregation")
	fmt.Println("• Use WaitGroup to coordinate")
	fmt.Println("• Collect results via channels")
}

// ============================================================================
// HTTP API INTEGRATION
// ============================================================================
// Demonstrates using concurrency in a web service

// App manages the RSS reader application
type App struct {
	store      *ArticleStore
	workerPool *WorkerPool
}

// NewApp creates the application
func NewApp() *App {
	store := NewArticleStore()
	pool := NewWorkerPool(5, store)

	return &App{
		store:      store,
		workerPool: pool,
	}
}

// articlesHandler returns recent articles as JSON
func (a *App) articlesHandler(w http.ResponseWriter, r *http.Request) {
	count := 10
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if n, err := strconv.Atoi(countStr); err == nil && n > 0 {
			count = n
		}
	}

	articles := a.store.GetRecent(count)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// syncHandler fetches feeds synchronously and returns results
func (a *App) syncHandler(w http.ResponseWriter, r *http.Request) {
	urls := []string{
		"https://nytimes.com/rss",
		"https://bbc.com/rss",
		"https://cnn.com/rss",
	}

	feeds, err := fanOutFanIn(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add to store
	for _, feed := range feeds {
		a.store.AddArticles(feed.Articles)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fetched":  len(feeds),
		"articles": a.store.Count(),
	})
}

// ============================================================================
// CONCURRENCY BEST PRACTICES
// ============================================================================

func demonstrateBestPractices() {
	fmt.Println("\n=== CONCURRENCY BEST PRACTICES ===")

	fmt.Println("\n✓ DO:")
	fmt.Println("• Keep it simple - don't overcomplicate")
	fmt.Println("• Prefer channels over shared memory")
	fmt.Println("• Use sync.Mutex for shared data access")
	fmt.Println("• Pass context for cancellation/timeouts")
	fmt.Println("• Limit goroutine count with worker pools")
	fmt.Println("• Close channels when done (by sender)")
	fmt.Println("• Handle errors explicitly")
	fmt.Println("• Test with 'go test -race'")
	fmt.Println("• Use WaitGroup/errgroup to track goroutines")

	fmt.Println("\n✗ DON'T:")
	fmt.Println("• Launch unbounded goroutines")
	fmt.Println("• Share data without synchronization")
	fmt.Println("• Close channels from receiver")
	fmt.Println("• Use time.Sleep for synchronization")
	fmt.Println("• Store contexts in structs")
	fmt.Println("• Ignore errors from goroutines")
	fmt.Println("• Use 'go' without tracking completion")

	fmt.Println("\n📚 WHEN TO USE CONCURRENCY:")
	fmt.Println("• Handling multiple network requests")
	fmt.Println("• I/O-bound operations (files, databases)")
	fmt.Println("• Data pipelines and streaming")
	fmt.Println("• Background tasks")
	fmt.Println("• Worker pools for job processing")
	fmt.Println("• Real-time data feeds")
	fmt.Println("• Parallel computations (CPU-bound)")

	fmt.Println("\n🚫 WHEN NOT TO USE CONCURRENCY:")
	fmt.Println("• Simple, single-task programs")
	fmt.Println("• Tasks requiring strict ordering")
	fmt.Println("• Lack of concurrency experience")
	fmt.Println("• Overhead > benefits")
	fmt.Println("• Debugging is primary concern")
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                               ║")
	fmt.Println("║           CHAPTER 9: CONCURRENCY IN GO                        ║")
	fmt.Println("║                                                               ║")
	fmt.Println("║  Building an RSS Reader with Goroutines and Channels          ║")
	fmt.Println("║                                                               ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	// Run all demonstrations
	demonstrateBasicGoroutines()
	demonstrateWaitGroup()
	demonstrateErrGroup()
	demonstrateMutex()
	demonstrateContext()
	demonstrateBasicChannels()
	demonstrateDirectionalChannels()
	demonstrateSelect()
	demonstrateWorkerPool()
	demonstrateFanOutFanIn()
	demonstrateBestPractices()

	fmt.Println("\n" + "=================================================================")
	fmt.Println("HTTP API EXAMPLE:")
	fmt.Println("=================================================================")
	fmt.Println()
	fmt.Println("To run the HTTP server, uncomment the section below.")
	fmt.Println("The server will demonstrate concurrency in a real web service:")
	fmt.Println()
	fmt.Println("• GET /articles?count=10  - Get recent articles")
	fmt.Println("• GET /sync               - Fetch feeds concurrently")
	fmt.Println()

	// Uncomment to run HTTP server:
	/*
		app := NewApp()

		// Seed with some initial data
		initialURLs := []string{
			"https://nytimes.com/rss",
			"https://bbc.com/rss",
			"https://cnn.com/rss",
		}
		for _, url := range initialURLs {
			app.workerPool.Submit(url)
		}

		http.HandleFunc("/articles", app.articlesHandler)
		http.HandleFunc("/sync", app.syncHandler)

		fmt.Println("Server running on http://localhost:8080")
		fmt.Println("Try: curl http://localhost:8080/articles?count=5")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	*/

	fmt.Println("\n✅ All demonstrations complete!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("• Go makes concurrency simple with goroutines and channels")
	fmt.Println("• Use sync primitives to coordinate and protect shared data")
	fmt.Println("• Context provides cancellation and timeout control")
	fmt.Println("• Channels enable safe communication between goroutines")
	fmt.Println("• Worker pools and fan-out/fan-in are essential patterns")
	fmt.Println("• Always track goroutines and handle errors properly")
	fmt.Println("• Test with 'go test -race' to catch race conditions")
}
