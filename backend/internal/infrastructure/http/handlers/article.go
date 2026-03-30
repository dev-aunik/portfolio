// Package handlers provides thin HTTP handler functions delegating to application services.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	articlesvc "github.com/aunik/portfolio/internal/application/article"
	"github.com/aunik/portfolio/internal/domain/article"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ArticleHandler holds the article service dependency.
type ArticleHandler struct {
	svc *articlesvc.Service
}

// NewArticleHandler creates a new ArticleHandler.
func NewArticleHandler(svc *articlesvc.Service) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// List handles GET /api/v1/articles
func (h *ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 20)

	result, err := h.svc.List(r.Context(), page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

// Get handles GET /api/v1/articles/{slug}
func (h *ArticleHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	a, err := h.svc.GetBySlug(r.Context(), slug)
	if errors.Is(err, article.ErrNotFound) {
		respondError(w, http.StatusNotFound, "article not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, a)
}

// Search handles GET /api/v1/articles/search?q=
func (h *ArticleHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 10)

	result, err := h.svc.Search(r.Context(), q, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

// ─── Admin handlers ──────────────────────────────────────────────────────────

type createArticleRequest struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Summary     string   `json:"summary"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	PublishedAt *string  `json:"published_at"` // RFC3339 string or null
}

// AdminCreate handles POST /api/v1/admin/articles
func (h *ArticleHandler) AdminCreate(w http.ResponseWriter, r *http.Request) {
	var req createArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	in := articlesvc.CreateInput{
		Title:   req.Title,
		Slug:    req.Slug,
		Summary: req.Summary,
		Content: req.Content,
		Tags:    req.Tags,
	}
	if req.PublishedAt != nil {
		in.PublishedAt = parseTimePtr(*req.PublishedAt)
	}

	a, err := h.svc.Create(r.Context(), in)
	if err != nil {
		if errors.Is(err, article.ErrSlugTaken) {
			respondError(w, http.StatusConflict, "slug is already taken")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, a)
}

// AdminUpdate handles PUT /api/v1/admin/articles/{id}
func (h *ArticleHandler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid article id")
		return
	}

	var req createArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	in := articlesvc.UpdateInput{
		ID:      id,
		Title:   req.Title,
		Slug:    req.Slug,
		Summary: req.Summary,
		Content: req.Content,
		Tags:    req.Tags,
	}
	if req.PublishedAt != nil {
		in.PublishedAt = parseTimePtr(*req.PublishedAt)
	}

	a, err := h.svc.Update(r.Context(), in)
	if errors.Is(err, article.ErrNotFound) {
		respondError(w, http.StatusNotFound, "article not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, a)
}

// AdminDelete handles DELETE /api/v1/admin/articles/{id}
func (h *ArticleHandler) AdminDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid article id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, article.ErrNotFound) {
			respondError(w, http.StatusNotFound, "article not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
