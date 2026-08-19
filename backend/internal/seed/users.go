package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"toefl-prep/backend/internal/auth"
)

// UserSeed is a dev bootstrap account. Passwords are the documented defaults
// from the README; rotate them in production via a password update path.
type UserSeed struct {
	Email    string
	Password string
	Role     string
}

func Users() []UserSeed {
	return []UserSeed{
		{Email: "student", Password: "123", Role: "student"},
		{Email: "admin", Password: "123", Role: "admin"},
	}
}

// EnsureUsers inserts the bootstrap accounts if their email is missing.
func EnsureUsers(ctx context.Context, pool *pgxpool.Pool) error {
	for _, u := range Users() {
		var exists bool
		if err := pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, u.Email).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		hash, err := auth.HashPassword(u.Password)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx,
			`INSERT INTO users (email, password_hash, role) VALUES ($1,$2,$3)`,
			u.Email, hash, u.Role); err != nil {
			return err
		}
	}
	return nil
}