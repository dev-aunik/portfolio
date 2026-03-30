// Package main is the entrypoint for the portfolio backend server.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	articleapp "github.com/aunik/portfolio/internal/application/article"
	contactapp "github.com/aunik/portfolio/internal/application/contact"
	"github.com/aunik/portfolio/internal/infrastructure/cache"
	httpinfra "github.com/aunik/portfolio/internal/infrastructure/http"
	"github.com/aunik/portfolio/internal/infrastructure/http/handlers"
	"github.com/aunik/portfolio/internal/infrastructure/messaging"
	"github.com/aunik/portfolio/internal/infrastructure/persistence"
	"github.com/aunik/portfolio/internal/infrastructure/search"
	"github.com/aunik/portfolio/pkg/config"
	"github.com/aunik/portfolio/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// ─── Config ────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// ─── Logger ────────────────────────────────────────────────────────────
	log := logger.New(cfg.Log.Level)
	log.Info().Str("env", cfg.App.Env).Msg("starting portfolio backend")

	// ─── Database ──────────────────────────────────────────────────────────
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("postgres ping failed")
	}
	log.Info().Str("host", cfg.DB.Host).Msg("postgres connected")

	// ─── Cache ─────────────────────────────────────────────────────────────
	redisCache, err := cache.New(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	defer redisCache.Close()
	log.Info().Str("host", cfg.Redis.Host).Msg("redis connected")

	// ─── Search ────────────────────────────────────────────────────────────
	searchEngine, err := search.New(
		cfg.Typesense.Host,
		cfg.Typesense.Port,
		cfg.Typesense.Protocol,
		cfg.Typesense.APIKey,
		cfg.Typesense.Collection,
	)
	if err != nil {
		log.Warn().Err(err).Msg("typesense unavailable — search disabled")
		searchEngine = nil
	} else {
		log.Info().Str("host", cfg.Typesense.Host).Msg("typesense connected")
	}

	// ─── Messaging (RabbitMQ) ──────────────────────────────────────────────
	rabbitBus, err := messaging.NewRabbitMQ(cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.QueueContact)
	if err != nil {
		log.Warn().Err(err).Msg("rabbitmq unavailable — contact queueing disabled")
		rabbitBus = nil
	} else {
		log.Info().Str("url", cfg.RabbitMQ.URL).Msg("rabbitmq connected")
		defer rabbitBus.Close()
	}

	// ─── Messaging (Kafka/Redpanda — audit log) ────────────────────────────
	kafkaBus, err := messaging.NewKafka(cfg.Kafka.Brokers, cfg.Kafka.TopicAudit)
	if err != nil {
		log.Warn().Err(err).Msg("kafka unavailable — audit log disabled")
		kafkaBus = nil
	} else {
		log.Info().Strs("brokers", cfg.Kafka.Brokers).Msg("kafka connected")
		defer kafkaBus.Close()
	}

	// ─── Repositories ─────────────────────────────────────────────────────
	articleRepo := persistence.NewArticleRepo(dbPool)
	contactRepo := persistence.NewContactRepo(dbPool)

	// ─── Application Services ─────────────────────────────────────────────
	articleSvc := articleapp.NewService(
		articleRepo,
		redisCache,
		searchEngine,
		kafkaBus,
		cfg.Redis.ArticlesTTL,
		cfg.Redis.ArticleTTL,
	)
	contactSvc := contactapp.NewService(contactRepo, rabbitBus)

	// ─── HTTP Handlers ─────────────────────────────────────────────────────
	articleH := handlers.NewArticleHandler(articleSvc)
	contactH := handlers.NewContactHandler(contactSvc)
	authH := handlers.NewAuthHandler(
		cfg.Admin.Email,
		cfg.Admin.Password,
		cfg.JWT.Secret,
		cfg.JWT.ExpiryMinutes,
	)

	// ─── Router ────────────────────────────────────────────────────────────
	router := httpinfra.NewRouter(httpinfra.Config{
		ArticleH:       articleH,
		ContactH:       contactH,
		AuthH:          authH,
		AllowedOrigins: cfg.App.AllowedOrigins,
		JWTSecret:      cfg.JWT.Secret,
		RatePublic:     cfg.RateLimit.Public,
		RateAdmin:      cfg.RateLimit.Admin,
		RateContact:    cfg.RateLimit.Contact,
	})

	// ─── HTTP Server ───────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// ─── Graceful Shutdown ─────────────────────────────────────────────────
	done := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Info().Msg("shutdown signal received")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("server shutdown error")
		}
		close(done)
	}()

	log.Info().Int("port", cfg.App.Port).Msg("server listening")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server error")
	}

	<-done
	log.Info().Msg("server stopped gracefully")
}
