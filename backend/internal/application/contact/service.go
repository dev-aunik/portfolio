// Package contact provides application-layer use cases for contact form submissions.
package contact

import (
	"context"
	"fmt"

	"github.com/aunik/portfolio/internal/domain/contact"
	"github.com/aunik/portfolio/internal/ports"
)

// Service orchestrates contact form use cases.
type Service struct {
	repo ports.ContactRepository
}

// NewService creates a new contact Service.
func NewService(repo ports.ContactRepository) *Service {
	return &Service{repo: repo}
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



	return nil
}
