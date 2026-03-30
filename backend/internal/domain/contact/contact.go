// Package contact defines the core Contact domain entity and related types.
package contact

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Status represents the lifecycle state of a contact submission.
type Status string

const (
	StatusPending   Status = "pending"
	StatusProcessed Status = "processed"
	StatusFailed    Status = "failed"
)

// ErrValidation is returned when contact form data is invalid.
var ErrValidation = errors.New("contact: validation failed")

// Contact is the core domain entity for a contact form submission.
type Contact struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Subject   string
	Message   string
	Status    Status
	CreatedAt time.Time
}

// Validate performs domain-level validation on the contact submission.
func (c *Contact) Validate() error {
	var errs []string

	if strings.TrimSpace(c.Name) == "" {
		errs = append(errs, "name is required")
	}
	if utf8.RuneCountInString(c.Name) > 100 {
		errs = append(errs, "name must be 100 characters or fewer")
	}

	if strings.TrimSpace(c.Email) == "" {
		errs = append(errs, "email is required")
	} else if _, err := mail.ParseAddress(c.Email); err != nil {
		errs = append(errs, "email is not valid")
	}

	if strings.TrimSpace(c.Message) == "" {
		errs = append(errs, "message is required")
	}
	if utf8.RuneCountInString(c.Message) > 5000 {
		errs = append(errs, "message must be 5000 characters or fewer")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrValidation, strings.Join(errs, "; "))
	}
	return nil
}

// CreateParams holds the raw input from a user contact form.
type CreateParams struct {
	Name    string
	Email   string
	Subject string
	Message string
}

// New creates a validated Contact entity from user input.
func New(p CreateParams) (*Contact, error) {
	c := &Contact{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(p.Name),
		Email:     strings.ToLower(strings.TrimSpace(p.Email)),
		Subject:   strings.TrimSpace(p.Subject),
		Message:   strings.TrimSpace(p.Message),
		Status:    StatusPending,
		CreatedAt: time.Now().UTC(),
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}
