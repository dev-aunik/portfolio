// Package config loads and exposes typed application configuration via Viper.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App       AppConfig
	DB        DBConfig
	Redis     RedisConfig
	Typesense TypesenseConfig
	RabbitMQ  RabbitMQConfig
	Kafka     KafkaConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Admin     AdminConfig
	SMTP      SMTPConfig
	Log       LogConfig
}

type AppConfig struct {
	Env            string
	Debug          bool
	Port           int
	AllowedOrigins []string
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	URL      string
}

type RedisConfig struct {
	Host        string
	Port        int
	Password    string
	DB          int
	URL         string
	ArticlesTTL time.Duration
	ArticleTTL  time.Duration
}

type TypesenseConfig struct {
	Host       string
	Port       int
	Protocol   string
	APIKey     string
	Collection string
}

type RabbitMQConfig struct {
	URL          string
	Exchange     string
	QueueContact string
}

type KafkaConfig struct {
	Brokers      []string
	TopicAudit   string
	TopicContact string
	GroupID      string
}

type JWTConfig struct {
	Secret             string
	ExpiryMinutes      int
	RefreshExpiryHours int
}

type RateLimitConfig struct {
	Public  int
	Admin   int
	Contact int
}

type AdminConfig struct {
	Email    string
	Password string
}

type SMTPConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	From        string
	NotifyEmail string
}

type LogConfig struct {
	Level string
}

// Load reads configuration from environment variables and .env file.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig() // OK if .env missing in production (env vars set externally)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_DEBUG", false)
	v.SetDefault("BACKEND_PORT", 8080)
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("CACHE_ARTICLES_TTL", 300)
	v.SetDefault("CACHE_ARTICLE_TTL", 600)
	v.SetDefault("TYPESENSE_PORT", 8108)
	v.SetDefault("TYPESENSE_PROTOCOL", "http")
	v.SetDefault("TYPESENSE_COLLECTION_ARTICLES", "articles")
	v.SetDefault("JWT_EXPIRY_MINUTES", 60)
	v.SetDefault("JWT_REFRESH_EXPIRY_HOURS", 168)
	v.SetDefault("RATE_LIMIT_PUBLIC", 100)
	v.SetDefault("RATE_LIMIT_ADMIN", 1000)
	v.SetDefault("RATE_LIMIT_CONTACT", 5)
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("KAFKA_TOPIC_AUDIT", "portfolio.audit")
	v.SetDefault("KAFKA_TOPIC_CONTACT", "portfolio.contact")
	v.SetDefault("KAFKA_GROUP_ID", "portfolio-consumer")
	v.SetDefault("SMTP_PORT", 587)

	originsRaw := v.GetString("ALLOWED_ORIGINS")
	var origins []string
	for _, o := range strings.Split(originsRaw, ",") {
		if t := strings.TrimSpace(o); t != "" {
			origins = append(origins, t)
		}
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:4200"}
	}

	brokersRaw := v.GetString("KAFKA_BROKERS")
	var brokers []string
	for _, b := range strings.Split(brokersRaw, ",") {
		if t := strings.TrimSpace(b); t != "" {
			brokers = append(brokers, t)
		}
	}
	if len(brokers) == 0 {
		brokers = []string{"localhost:9092"}
	}

	articlesTTL := time.Duration(v.GetInt("CACHE_ARTICLES_TTL")) * time.Second
	articleTTL := time.Duration(v.GetInt("CACHE_ARTICLE_TTL")) * time.Second

	cfg := &Config{
		App: AppConfig{
			Env:            v.GetString("APP_ENV"),
			Debug:          v.GetBool("APP_DEBUG"),
			Port:           v.GetInt("BACKEND_PORT"),
			AllowedOrigins: origins,
		},
		DB: DBConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt("DB_PORT"),
			Name:     v.GetString("DB_NAME"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			SSLMode:  v.GetString("DB_SSLMODE"),
			URL:      v.GetString("DATABASE_URL"),
		},
		Redis: RedisConfig{
			Host:        v.GetString("REDIS_HOST"),
			Port:        v.GetInt("REDIS_PORT"),
			Password:    v.GetString("REDIS_PASSWORD"),
			DB:          v.GetInt("REDIS_DB"),
			URL:         v.GetString("REDIS_URL"),
			ArticlesTTL: articlesTTL,
			ArticleTTL:  articleTTL,
		},
		Typesense: TypesenseConfig{
			Host:       v.GetString("TYPESENSE_HOST"),
			Port:       v.GetInt("TYPESENSE_PORT"),
			Protocol:   v.GetString("TYPESENSE_PROTOCOL"),
			APIKey:     v.GetString("TYPESENSE_API_KEY"),
			Collection: v.GetString("TYPESENSE_COLLECTION_ARTICLES"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:          v.GetString("RABBITMQ_URL"),
			Exchange:     v.GetString("RABBITMQ_EXCHANGE"),
			QueueContact: v.GetString("RABBITMQ_QUEUE_CONTACT"),
		},
		Kafka: KafkaConfig{
			Brokers:      brokers,
			TopicAudit:   v.GetString("KAFKA_TOPIC_AUDIT"),
			TopicContact: v.GetString("KAFKA_TOPIC_CONTACT"),
			GroupID:      v.GetString("KAFKA_GROUP_ID"),
		},
		JWT: JWTConfig{
			Secret:             v.GetString("JWT_SECRET"),
			ExpiryMinutes:      v.GetInt("JWT_EXPIRY_MINUTES"),
			RefreshExpiryHours: v.GetInt("JWT_REFRESH_EXPIRY_HOURS"),
		},
		RateLimit: RateLimitConfig{
			Public:  v.GetInt("RATE_LIMIT_PUBLIC"),
			Admin:   v.GetInt("RATE_LIMIT_ADMIN"),
			Contact: v.GetInt("RATE_LIMIT_CONTACT"),
		},
		Admin: AdminConfig{
			Email:    v.GetString("ADMIN_EMAIL"),
			Password: v.GetString("ADMIN_PASSWORD"),
		},
		SMTP: SMTPConfig{
			Host:        v.GetString("SMTP_HOST"),
			Port:        v.GetInt("SMTP_PORT"),
			User:        v.GetString("SMTP_USER"),
			Password:    v.GetString("SMTP_PASSWORD"),
			From:        v.GetString("SMTP_FROM"),
			NotifyEmail: v.GetString("NOTIFY_EMAIL"),
		},
		Log: LogConfig{
			Level: v.GetString("LOG_LEVEL"),
		},
	}

	if cfg.DB.URL == "" && cfg.DB.Host != "" {
		cfg.DB.URL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)
	}

	return cfg, nil
}
