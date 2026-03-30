// Package persistence provides a PostgreSQL-backed implementation of repository ports.
package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aunik/portfolio/internal/domain/article"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ArticleRepo implements ports.ArticleRepository using PostgreSQL via pgx.
type ArticleRepo struct {
	db *pgxpool.Pool
}

// NewArticleRepo creates a new ArticleRepo backed by the given pool.
func NewArticleRepo(db *pgxpool.Pool) *ArticleRepo {
	return &ArticleRepo{db: db}
}

// FindAll returns a paginated list of published articles, newest first.
func (r *ArticleRepo) FindAll(ctx context.Context, page, limit int) ([]*article.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Count total published
	var total int64
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM articles WHERE published_at IS NOT NULL AND published_at <= NOW()`).
		Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("articles: count: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, title, slug, summary, content, tags, published_at, created_at, updated_at
		FROM articles
		WHERE published_at IS NOT NULL AND published_at <= NOW()
		ORDER BY published_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("articles: find all: %w", err)
	}
	defer rows.Close()

	articles, err := scanArticles(rows)
	if err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// FindBySlug retrieves a single published article by slug.
func (r *ArticleRepo) FindBySlug(ctx context.Context, slug string) (*article.Article, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, title, slug, summary, content, tags, published_at, created_at, updated_at
		FROM articles
		WHERE slug = $1 AND published_at IS NOT NULL AND published_at <= NOW()
	`, slug)

	a, err := scanArticle(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, article.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("articles: find by slug %q: %w", slug, err)
	}
	return a, nil
}

// FindByID retrieves any article by ID (including unpublished, for admin).
func (r *ArticleRepo) FindByID(ctx context.Context, id uuid.UUID) (*article.Article, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, title, slug, summary, content, tags, published_at, created_at, updated_at
		FROM articles WHERE id = $1
	`, id)

	a, err := scanArticle(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, article.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("articles: find by id %q: %w", id, err)
	}
	return a, nil
}

// Create inserts a new article.
func (r *ArticleRepo) Create(ctx context.Context, a *article.Article) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO articles (id, title, slug, summary, content, tags, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, a.ID, a.Title, a.Slug, a.Summary, a.Content, a.Tags, a.PublishedAt, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique") && strings.Contains(err.Error(), "slug") {
			return article.ErrSlugTaken
		}
		return fmt.Errorf("articles: create: %w", err)
	}
	return nil
}

// Update saves changes to an existing article.
func (r *ArticleRepo) Update(ctx context.Context, a *article.Article) error {
	a.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx, `
		UPDATE articles
		SET title=$1, slug=$2, summary=$3, content=$4, tags=$5, published_at=$6, updated_at=$7
		WHERE id=$8
	`, a.Title, a.Slug, a.Summary, a.Content, a.Tags, a.PublishedAt, a.UpdatedAt, a.ID)
	if err != nil {
		if strings.Contains(err.Error(), "unique") && strings.Contains(err.Error(), "slug") {
			return article.ErrSlugTaken
		}
		return fmt.Errorf("articles: update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return article.ErrNotFound
	}
	return nil
}

// Delete removes an article by ID.
func (r *ArticleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM articles WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("articles: delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return article.ErrNotFound
	}
	return nil
}

// ─── helpers ────────────────────────────────────────────────────────────────

func scanArticle(row pgx.Row) (*article.Article, error) {
	var a article.Article
	var tags []string
	err := row.Scan(
		&a.ID, &a.Title, &a.Slug, &a.Summary, &a.Content,
		&tags, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	a.Tags = tags
	return &a, nil
}

func scanArticles(rows pgx.Rows) ([]*article.Article, error) {
	var result []*article.Article
	for rows.Next() {
		var a article.Article
		var tags []string
		if err := rows.Scan(
			&a.ID, &a.Title, &a.Slug, &a.Summary, &a.Content,
			&tags, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("articles: scan row: %w", err)
		}
		a.Tags = tags
		result = append(result, &a)
	}
	return result, rows.Err()
}
