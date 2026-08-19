package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	httppkg "toefl-prep/backend/internal/http"
)

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Service struct {
	pool *pgxpool.Pool
	ttl  time.Duration
}

func New(pool *pgxpool.Pool, ttl time.Duration) *Service {
	return &Service{pool: pool, ttl: ttl}
}

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Login verifies credentials and returns a fresh session token.
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	var (
		userID int64
		hash   string
		role   string
	)
	err := s.pool.QueryRow(ctx,
		`SELECT id, password_hash, role FROM users WHERE email = $1`, email,
	).Scan(&userID, &hash, &role)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", httppkg.ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return "", httppkg.ErrUnauthorized
	}

	token, err := newToken()
	if err != nil {
		return "", err
	}
	_, err = s.pool.Exec(ctx,
		`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, time.Now().Add(s.ttl),
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Verify(ctx context.Context, token string) (*User, error) {
	if token == "" {
		return nil, httppkg.ErrUnauthorized
	}
	var (
		userID    int64
		email     string
		role      string
		expiresAt time.Time
	)
	err := s.pool.QueryRow(ctx,
		`SELECT s.user_id, u.email, u.role, s.expires_at
		 FROM sessions s JOIN users u ON u.id = s.user_id
		 WHERE s.token = $1`, token,
	).Scan(&userID, &email, &role, &expiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, httppkg.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	if time.Now().After(expiresAt) {
		_, _ = s.pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
		return nil, httppkg.ErrUnauthorized
	}
	return &User{ID: userID, Email: email, Role: role}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}