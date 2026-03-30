// Package handlers — contact form handler.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	contactsvc "github.com/aunik/portfolio/internal/application/contact"
	"github.com/aunik/portfolio/internal/domain/contact"
)

// ContactHandler holds the contact service dependency.
type ContactHandler struct {
	svc *contactsvc.Service
}

// NewContactHandler creates a new ContactHandler.
func NewContactHandler(svc *contactsvc.Service) *ContactHandler {
	return &ContactHandler{svc: svc}
}

type submitContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Submit handles POST /api/v1/contact
func (h *ContactHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req submitContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.svc.Submit(r.Context(), contactsvc.SubmitInput{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	})
	if errors.Is(err, contact.ErrValidation) {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to process submission")
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]string{
		"message": "Thank you! Your message has been received.",
	})
}
