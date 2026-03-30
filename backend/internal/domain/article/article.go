// Package article defines the core Article domain entity and related value objects.
package article

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ErrNotFound is returned when an article cannot be found.
var ErrNotFound = errors.New("article: not found")

// ErrSlugTaken is returned when an article slug is already in use.
var ErrSlugTaken = errors.New("article: slug already taken")

// Article is the core domain entity for a blog post.
type Article struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content"`
	Tags        []string   `json:"tags"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// IsPublished reports whether the article has been published.
func (a *Article) IsPublished() bool {
	return a.PublishedAt != nil && !a.PublishedAt.IsZero() && a.PublishedAt.Before(time.Now())
}

// ReadingMinutes estimates the reading time in minutes based on 200 wpm.
func (a *Article) ReadingMinutes() int {
	words := len(strings.Fields(a.Content))
	mins := words / 200
	if mins < 1 {
		return 1
	}
	return mins
}

// Validate performs basic domain validation on the article.
func (a *Article) Validate() error {
	if strings.TrimSpace(a.Title) == "" {
		return errors.New("article: title is required")
	}
	if utf8.RuneCountInString(a.Title) > 200 {
		return errors.New("article: title exceeds 200 characters")
	}
	if strings.TrimSpace(a.Slug) == "" {
		return errors.New("article: slug is required")
	}
	if strings.TrimSpace(a.Summary) == "" {
		return errors.New("article: summary is required")
	}
	if strings.TrimSpace(a.Content) == "" {
		return errors.New("article: content is required")
	}
	return nil
}

// CreateParams holds the input for creating a new article.
type CreateParams struct {
	Title       string
	Slug        string
	Summary     string
	Content     string
	Tags        []string
	PublishedAt *time.Time
}

// New creates a new Article from CreateParams.
func New(p CreateParams) (*Article, error) {
	a := &Article{
		ID:          uuid.New(),
		Title:       strings.TrimSpace(p.Title),
		Slug:        strings.TrimSpace(p.Slug),
		Summary:     strings.TrimSpace(p.Summary),
		Content:     p.Content,
		Tags:        p.Tags,
		PublishedAt: p.PublishedAt,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if a.Tags == nil {
		a.Tags = []string{}
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return a, nil
}
