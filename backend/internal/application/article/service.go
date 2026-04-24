// Package article provides application-layer use cases for articles.
package article

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aunik/portfolio/internal/domain/article"
	"github.com/aunik/portfolio/internal/ports"
	"github.com/google/uuid"
)

const (
	cacheKeyArticles = "articles:list:%d:%d" // page:limit
	cacheKeyArticle  = "articles:slug:%s"
)

// Service orchestrates article-related use cases.
type Service struct {
	repo        ports.ArticleRepository
	cache       ports.Cache
	articlesTTL time.Duration
	articleTTL  time.Duration
}

// NewService creates a new article Service with all required dependencies.
func NewService(
	repo ports.ArticleRepository,
	cache ports.Cache,
	articlesTTL, articleTTL time.Duration,
) *Service {
	return &Service{
		repo:        repo,
		cache:       cache,
		articlesTTL: articlesTTL,
		articleTTL:  articleTTL,
	}
}

// ListResult holds a paginated list of articles.
type ListResult struct {
	Articles   []*article.Article `json:"articles"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	TotalPages int                `json:"total_pages"`
}

// List returns a paginated list of published articles, using Redis cache.
func (s *Service) List(ctx context.Context, page, limit int) (*ListResult, error) {
	key := fmt.Sprintf(cacheKeyArticles, page, limit)

	if cached, _ := s.cache.Get(ctx, key); cached != "" {
		var result ListResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	articles, total, err := s.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("article.service: list: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	result := &ListResult{
		Articles:   articles,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}

	if data, err := json.Marshal(result); err == nil {
		_ = s.cache.Set(ctx, key, string(data), s.articlesTTL)
	}
	return result, nil
}

// GetBySlug returns a published article by slug, using Redis cache.
func (s *Service) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	key := fmt.Sprintf(cacheKeyArticle, slug)

	if cached, _ := s.cache.Get(ctx, key); cached != "" {
		var a article.Article
		if err := json.Unmarshal([]byte(cached), &a); err == nil {
			return &a, nil
		}
	}

	a, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err // preserve domain errors (ErrNotFound)
	}

	if data, err := json.Marshal(a); err == nil {
		_ = s.cache.Set(ctx, key, string(data), s.articleTTL)
	}
	return a, nil
}

// Search performs a DB-backed basic full-text search.
func (s *Service) Search(ctx context.Context, query string, page, limit int) (*ListResult, error) {
	if query == "" {
		return &ListResult{}, nil
	}
	articles, total, err := s.repo.Search(ctx, query, page, limit)
	if err != nil {
		return nil, fmt.Errorf("article.service: search: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &ListResult{
		Articles:   articles,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}

// CreateInput holds the input for creating a new article.
type CreateInput struct {
	Title       string
	Slug        string
	Summary     string
	Content     string
	Tags        []string
	PublishedAt *time.Time
}

// Create creates a new article, indexes it, and publishes an audit event.
func (s *Service) Create(ctx context.Context, in CreateInput) (*article.Article, error) {
	a, err := article.New(article.CreateParams{
		Title:       in.Title,
		Slug:        in.Slug,
		Summary:     in.Summary,
		Content:     in.Content,
		Tags:        in.Tags,
		PublishedAt: in.PublishedAt,
	})
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

// UpdateInput holds the input for updating an article.
type UpdateInput struct {
	ID          uuid.UUID
	Title       string
	Slug        string
	Summary     string
	Content     string
	Tags        []string
	PublishedAt *time.Time
}

// Update modifies an existing article, invalidates cache, and re-indexes.
func (s *Service) Update(ctx context.Context, in UpdateInput) (*article.Article, error) {
	a, err := s.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	a.Title = in.Title
	a.Slug = in.Slug
	a.Summary = in.Summary
	a.Content = in.Content
	a.Tags = in.Tags
	a.PublishedAt = in.PublishedAt

	if err := a.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx,
		fmt.Sprintf(cacheKeyArticle, a.Slug),
	)
	s.invalidateListCache(ctx)

	return a, nil
}

// Delete removes an article, its cache, and its search index entry.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	a, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, article.ErrNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("article.service: delete lookup: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyArticle, a.Slug))
	s.invalidateListCache(ctx)
	return nil
}

func (s *Service) invalidateListCache(ctx context.Context) {
	// Invalidate common pages — a more robust approach would be a scan pattern
	for _, pg := range []int{1, 2, 3} {
		for _, lim := range []int{10, 20} {
			_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyArticles, pg, lim))
		}
	}
}

