// Package http provides the Chi router setup for the portfolio API.
package http

import (
	"net/http"
	"time"

	"github.com/aunik/portfolio/internal/infrastructure/http/handlers"
	mw "github.com/aunik/portfolio/internal/infrastructure/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config holds the dependencies needed to build the router.
type Config struct {
	ArticleH       *handlers.ArticleHandler
	ContactH       *handlers.ContactHandler
	AuthH          *handlers.AuthHandler
	AllowedOrigins []string
	JWTSecret      string
	RatePublic     int // requests per minute
	RateAdmin      int
	RateContact    int
}

// NewRouter builds and returns the fully-configured Chi router.
func NewRouter(cfg Config) http.Handler {
	r := chi.NewRouter()

	// ─── Global middleware (must all come before any routes) ────────────────
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(mw.SecurityHeaders)
	r.Use(mw.PrometheusMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// ─── Metrics ────────────────────────────────────────────────────────────
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// ─── Public API ─────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Article endpoints (cached)
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(cfg.RatePublic, time.Minute))

			r.Get("/articles", cfg.ArticleH.List)
			r.Get("/articles/search", cfg.ArticleH.Search) // before /{slug}
			r.Get("/articles/{slug}", cfg.ArticleH.Get)
		})

		// Contact endpoint (stricter rate limit)
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(cfg.RateContact, time.Minute))
			r.Post("/contact", cfg.ContactH.Submit)
		})

		// Admin – auth
		r.Post("/admin/login", cfg.AuthH.Login)

		// Admin – protected article management
		r.Group(func(r chi.Router) {
			r.Use(mw.Authenticate(cfg.JWTSecret))
			r.Use(httprate.LimitByIP(cfg.RateAdmin, time.Minute))

			r.Post("/admin/articles", cfg.ArticleH.AdminCreate)
			r.Put("/admin/articles/{id}", cfg.ArticleH.AdminUpdate)
			r.Delete("/admin/articles/{id}", cfg.ArticleH.AdminDelete)
		})
	})

	// ─── 404 catch-all ─────────────────────────────────────────────────────
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"route not found"}`))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error":"method not allowed"}`))
	})

	return r
}
