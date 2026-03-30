// Package persistence — contact repository for PostgreSQL.
package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ContactRepo implements ports.ContactRepository.
type ContactRepo struct {
	db *pgxpool.Pool
}

// NewContactRepo creates a new ContactRepo.
func NewContactRepo(db *pgxpool.Pool) *ContactRepo {
	return &ContactRepo{db: db}
}

// Create saves a new contact submission to the database.
func (r *ContactRepo) Create(ctx context.Context, name, email, subject, message string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO contacts (name, email, subject, message, status)
		VALUES ($1, $2, $3, $4, 'pending')
	`, name, email, subject, message)
	if err != nil {
		return fmt.Errorf("contacts: create: %w", err)
	}
	return nil
}
