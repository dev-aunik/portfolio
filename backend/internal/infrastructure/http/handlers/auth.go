// Package handlers — admin authentication handler.
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	mw "github.com/aunik/portfolio/internal/infrastructure/http/middleware"
	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler handles admin authentication.
type AuthHandler struct {
	adminEmail    string
	adminPassword string
	jwtSecret     string
	expiryMinutes int
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(email, password, jwtSecret string, expiryMinutes int) *AuthHandler {
	return &AuthHandler{
		adminEmail:    email,
		adminPassword: password,
		jwtSecret:     jwtSecret,
		expiryMinutes: expiryMinutes,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Login handles POST /api/v1/admin/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Constant-time comparison would be ideal in production; for simplicity we check directly.
	if req.Email != h.adminEmail || req.Password != h.adminPassword {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	expiresAt := time.Now().UTC().Add(time.Duration(h.expiryMinutes) * time.Minute)
	claims := &mw.Claims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   req.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "token generation failed")
		return
	}

	respondJSON(w, http.StatusOK, loginResponse{
		Token:     tokenStr,
		ExpiresAt: expiresAt,
	})
}
