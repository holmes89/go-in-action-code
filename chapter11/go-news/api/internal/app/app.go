package app

import (
	"context"
	"fmt"

	"github.com/goinaction/go-news/api/internal/feed"
	"github.com/goinaction/go-news/api/internal/store"
	"github.com/goinaction/go-news/api/internal/worker"
)

// App is the main application struct that coordinates all components.
type App struct {
	store      *store.ArticleStore
	feedRead   chan string
	errChan    chan error
	workerPool *worker.WorkerPool
	fetcher    feed.Fetcher
}

// NewApp creates and initializes a new App instance with all necessary components.
func NewApp() *App {
	feedReader := make(chan string, 5)
	errChan := make(chan error, 5)
	articleStore := store.NewArticleStore()
	feedFetcher := feed.NewGoFeedAdapter()
	
	app := &App{
		store:    articleStore,
		feedRead: feedReader,
		errChan:  errChan,
		fetcher:  feedFetcher,
	}
	
	app.workerPool = worker.NewWorkerPool(app, 3)
	go app.run()
	
	return app
}

// run continuously monitors channels for incoming feed URLs and errors.
// It uses a select statement to handle multiple channel operations.
func (a *App) run() {
	for {
		select {
		case feed := <-a.feedRead:
			a.workerPool.Submit(feed)
		case err := <-a.errChan:
			fmt.Println("error:", err)
		default:
			// Non-blocking: no operation if no messages
		}
	}
}

// ProcessFeed retrieves and processes an RSS feed from the given URL.
// It implements the FeedProcessor interface used by workers.
func (a *App) ProcessFeed(ctx context.Context, url string) ([]*feed.Article, error) {
	fmt.Println("fetching feed: ", url)
	
	feedData, err := a.fetcher.FetchFeed(ctx, url)
	if err != nil {
		fmt.Println("errored: ", url)
		a.errChan <- err
		return nil, err
	}
	fmt.Println("fetched: ", url)
	
	// Create a channel to stream articles to the save goroutine
	out := make(chan *feed.Article, 10)
	go a.saveFeed(out)
	defer close(out)
	
	for _, article := range feedData.Articles {
		fmt.Println("ingesting article: ", article.Link)
		out <- article
	}
	fmt.Println("processed: ", url)
	
	return feedData.Articles, nil
}

// saveFeed receives articles from a read-only channel and adds them to the store.
func (a *App) saveFeed(articles <-chan *feed.Article) {
	fmt.Println("adding articles...")
	for article := range articles {
		fmt.Println("adding article: ", article.Link)
		a.store.AddArticle(article)
		fmt.Println("article added: ", article.Link)
	}
	fmt.Println("all articles added")
}

// Store returns the application's article store for read operations.
func (a *App) Store() *store.ArticleStore {
	return a.store
}

// SubmitFeed queues a feed URL for processing by the worker pool.
func (a *App) SubmitFeed(url string) {
	a.feedRead <- url
}

// SubmitError sends an error to the error channel.
func (a *App) SubmitError(err error) {
	a.errChan <- err
}

// FetchFeedSync fetches a feed synchronously and returns the articles.
// This is used by HTTP handlers that need immediate results.
func (a *App) FetchFeedSync(ctx context.Context, url string) ([]*feed.Article, error) {
	return a.ProcessFeed(ctx, url)
}

// Close gracefully shuts down the application by closing channels
// and stopping the worker pool.
func (a *App) Close() {
	close(a.feedRead)
	close(a.errChan)
	a.workerPool.Stop()
}
