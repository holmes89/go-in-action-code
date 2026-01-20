package worker

import (
	"context"
	"fmt"
	"time"

	feed "github.com/goinaction/go-news/api/internal"
)

// FeedProcessor defines the interface for processing fetched feeds.
// This allows the worker to be decoupled from the specific implementation.
type FeedProcessor interface {
	ProcessFeed(ctx context.Context, url string) ([]*feed.Article, error)
}

// Worker processes incoming feed URLs from a channel.
type Worker struct {
	processor     FeedProcessor
	incomingFeeds <-chan string
}

// run continuously processes feed URLs from the incomingFeeds channel.
// It creates a context with timeout for each feed fetch operation.
func (w *Worker) run() {
	for url := range w.incomingFeeds {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		w.processor.ProcessFeed(ctx, url)
		cancel()
	}
}

// NewWorker creates a new Worker and starts it in a goroutine.
func NewWorker(processor FeedProcessor, incomingFeeds <-chan string) *Worker {
	w := &Worker{
		processor:     processor,
		incomingFeeds: incomingFeeds,
	}
	go w.run()
	return w
}

// WorkerPool manages a pool of workers that process feed URLs concurrently.
type WorkerPool struct {
	workers       []*Worker
	incomingFeeds chan string
}

// NewWorkerPool creates a pool of workers to process feed URLs.
// count specifies the number of worker goroutines to create.
func NewWorkerPool(processor FeedProcessor, count int) *WorkerPool {
	incomingFeeds := make(chan string, 100)
	wp := &WorkerPool{
		incomingFeeds: incomingFeeds,
	}
	
	for i := 0; i < count; i++ {
		w := NewWorker(processor, incomingFeeds)
		wp.workers = append(wp.workers, w)
	}
	
	return wp
}

// Submit adds a feed URL to the worker pool's queue for processing.
func (wp *WorkerPool) Submit(url string) {
	wp.incomingFeeds <- url
}

// Stop closes the worker pool's input channel, causing all workers to
// finish processing their current tasks and exit.
func (wp *WorkerPool) Stop() {
	close(wp.incomingFeeds)
}

// PrintStatus outputs the current status of the worker pool.
func (wp *WorkerPool) PrintStatus() {
	fmt.Printf("Worker pool: %d workers, %d pending feeds\n", 
		len(wp.workers), len(wp.incomingFeeds))
}
