// Package contact provides application-layer use cases for contact form submissions.
package contact

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aunik/portfolio/internal/domain/contact"
	"github.com/aunik/portfolio/internal/ports"
)

// Service orchestrates contact form use cases.
type Service struct {
	repo ports.ContactRepository
	bus  ports.MessageBus
}

// NewService creates a new contact Service.
func NewService(repo ports.ContactRepository, bus ports.MessageBus) *Service {
	return &Service{repo: repo, bus: bus}
}

// SubmitInput holds the raw form input from the HTTP layer.
type SubmitInput struct {
	Name    string
	Email   string
	Subject string
	Message string
}

// Submit validates, persists, and enqueues a contact form submission.
func (s *Service) Submit(ctx context.Context, in SubmitInput) error {
	// Domain validation
	c, err := contact.New(contact.CreateParams{
		Name:    in.Name,
		Email:   in.Email,
		Subject: in.Subject,
		Message: in.Message,
	})
	if err != nil {
		return err // domain validation error, presented to caller
	}

	// Persist to PostgreSQL (best-effort — don't block user if DB is slow)
	if err := s.repo.Create(ctx, c.Name, c.Email, c.Subject, c.Message); err != nil {
		// Log but don't fail — the message bus is the primary delivery mechanism
		fmt.Printf("contact.service: persist failed: %v\n", err)
	}

	// Publish to RabbitMQ for async processing (email notification, CRM, etc.)
	if s.bus != nil {
		payload, _ := json.Marshal(map[string]string{
			"id":      c.ID.String(),
			"name":    c.Name,
			"email":   c.Email,
			"subject": c.Subject,
			"message": c.Message,
			"ts":      time.Now().UTC().Format(time.RFC3339),
		})
		if err := s.bus.Publish(ctx, "contact.new", payload); err != nil {
			// Log failure but still return success — submission is already saved
			fmt.Printf("contact.service: enqueue failed: %v\n", err)
		}
	}

	return nil
}
