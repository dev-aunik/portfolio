// Package ports defines the interfaces (ports) for the application's external dependencies.
// Implementations live in internal/infrastructure.
package ports

import (
	"context"
	"time"

	"github.com/aunik/portfolio/internal/domain/article"
	"github.com/google/uuid"
)

// ArticleRepository defines the persistence contract for articles.
type ArticleRepository interface {
	// FindAll returns a paginated list of published articles along with the total count.
	FindAll(ctx context.Context, page, limit int) ([]*article.Article, int64, error)

	// FindBySlug retrieves a single article by its URL slug.
	FindBySlug(ctx context.Context, slug string) (*article.Article, error)

	// FindByID retrieves a single article by UUID (admin use).
	FindByID(ctx context.Context, id uuid.UUID) (*article.Article, error)

	// Create persists a new article to the store.
	Create(ctx context.Context, a *article.Article) error

	// Update persists changes to an existing article.
	Update(ctx context.Context, a *article.Article) error

	// Delete removes an article from the store.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Cache defines the caching contract (backed by Redis).
type Cache interface {
	// Get retrieves a cached value by key. Returns ("", nil) if key not found.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a key/value pair with an optional TTL.
	Set(ctx context.Context, key, value string, ttl time.Duration) error

	// Delete removes one or more keys from the cache.
	Delete(ctx context.Context, keys ...string) error

	// Exists reports whether a key exists in the cache.
	Exists(ctx context.Context, key string) (bool, error)
}

// SearchEngine defines the full-text search contract (backed by Typesense).
type SearchEngine interface {
	// IndexArticle adds or updates an article document in the search index.
	IndexArticle(ctx context.Context, a *article.Article) error

	// SearchArticles performs a full-text search and returns matching articles.
	SearchArticles(ctx context.Context, query string, page, limit int) (*SearchResult, error)

	// DeleteArticle removes an article document from the search index.
	DeleteArticle(ctx context.Context, id string) error
}

// SearchResult holds paginated search results.
type SearchResult struct {
	Hits       []*article.Article `json:"hits"`
	TotalHits  int64              `json:"total_hits"`
	Page       int                `json:"page"`
	TotalPages int                `json:"total_pages"`
}

// MessageBus defines the async messaging contract.
// Backed by RabbitMQ for contact submissions, Kafka for audit events.
type MessageBus interface {
	// Publish sends a message payload to the given topic/routing key.
	Publish(ctx context.Context, topic string, payload []byte) error

	// Close gracefully shuts down the connection.
	Close() error
}

// ContactRepository defines the persistence contract for contact submissions.
type ContactRepository interface {
	// Create saves a contact submission.
	Create(ctx context.Context, name, email, subject, message string) error
}
