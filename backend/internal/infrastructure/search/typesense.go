// Package search provides a Typesense-backed implementation of ports.SearchEngine.
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aunik/portfolio/internal/domain/article"
	"github.com/aunik/portfolio/internal/ports"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

const schemaVersion = "1"

// TypesenseEngine implements ports.SearchEngine.
type TypesenseEngine struct {
	client     *typesense.Client
	collection string
}

// articleDoc is the Typesense document shape.
type articleDoc struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Summary     string   `json:"summary"`
	Tags        []string `json:"tags"`
	PublishedAt int64    `json:"published_at"` // unix timestamp
}

// New creates a TypesenseEngine and ensures the collection schema exists.
func New(host string, port int, protocol, apiKey, collection string) (*TypesenseEngine, error) {
	serverURL := fmt.Sprintf("%s://%s:%d", protocol, host, port)
	client := typesense.NewClient(
		typesense.WithServer(serverURL),
		typesense.WithAPIKey(apiKey),
		typesense.WithConnectionTimeout(5*time.Second),
	)

	eng := &TypesenseEngine{client: client, collection: collection}
	if err := eng.ensureCollection(context.Background()); err != nil {
		return nil, fmt.Errorf("search: ensure collection: %w", err)
	}
	return eng, nil
}

func (e *TypesenseEngine) ensureCollection(ctx context.Context) error {
	_, err := e.client.Collection(e.collection).Retrieve(ctx)
	if err == nil {
		return nil // already exists
	}

	schema := &api.CollectionSchema{
		Name: e.collection,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "slug", Type: "string", Index: pointer.True()},
			{Name: "summary", Type: "string"},
			{Name: "tags", Type: "string[]", Facet: pointer.True()},
			{Name: "published_at", Type: "int64", Sort: pointer.True()},
		},
		DefaultSortingField: pointer.String("published_at"),
	}

	_, err = e.client.Collections().Create(ctx, schema)
	return err
}

// IndexArticle upserts an article document.
func (e *TypesenseEngine) IndexArticle(ctx context.Context, a *article.Article) error {
	doc := articleDoc{
		ID:      a.ID.String(),
		Title:   a.Title,
		Slug:    a.Slug,
		Summary: a.Summary,
		Tags:    a.Tags,
	}
	if a.PublishedAt != nil {
		doc.PublishedAt = a.PublishedAt.Unix()
	}

	action := "upsert"
	_, err := e.client.Collection(e.collection).Documents().Import(ctx, []interface{}{doc}, &api.ImportDocumentsParams{
		Action: &action,
	})
	if err != nil {
		return fmt.Errorf("search: index article %q: %w", a.ID, err)
	}
	return nil
}

// SearchArticles performs a full-text search and returns matching articles.
func (e *TypesenseEngine) SearchArticles(ctx context.Context, query string, page, limit int) (*ports.SearchResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	params := &api.SearchCollectionParams{
		Q:       query,
		QueryBy: "title,summary,tags",
		Page:    pointer.Int(page),
		PerPage: pointer.Int(limit),
		SortBy:  pointer.String("published_at:desc"),
	}

	res, err := e.client.Collection(e.collection).Documents().Search(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search: query %q: %w", query, err)
	}

	result := &ports.SearchResult{
		Page: page,
	}
	if res.Found != nil {
		result.TotalHits = int64(*res.Found)
		totalPages := (int(*res.Found) + limit - 1) / limit
		result.TotalPages = totalPages
	}

	for _, hit := range *res.Hits {
		raw, _ := json.Marshal(hit.Document)
		var doc articleDoc
		if err := json.Unmarshal(raw, &doc); err != nil {
			continue
		}
		var publishedAt *time.Time
		if doc.PublishedAt > 0 {
			t := time.Unix(doc.PublishedAt, 0).UTC()
			publishedAt = &t
		}
		result.Hits = append(result.Hits, &article.Article{
			Title:       doc.Title,
			Slug:        doc.Slug,
			Summary:     doc.Summary,
			Tags:        doc.Tags,
			PublishedAt: publishedAt,
		})
	}
	return result, nil
}

// DeleteArticle removes an article document from the index.
func (e *TypesenseEngine) DeleteArticle(ctx context.Context, id string) error {
	_, err := e.client.Collection(e.collection).Document(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("search: delete article %q: %w", id, err)
	}
	return nil
}
