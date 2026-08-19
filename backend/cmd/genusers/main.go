// Command genusers creates N student accounts for load testing, since each
// user may hold only one active attempt per exam template (FR-4.10) and k6
// VUs must therefore not share the default seeded student account.
//
// Usage:
//
//	DATABASE_URL=postgres://toefl:toefl@localhost:5432/toefl_test \
//	  go run ./cmd/genusers -n 20
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"toefl-prep/backend/internal/auth"
	"toefl-prep/backend/internal/database"
)

func main() {
	n := flag.Int("n", 10, "number of student users to create")
	password := flag.String("password", "123", "password for every created user")
	flag.Parse()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is required")
	}
	ctx := context.Background()
	pool, err := database.New(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	hash, err := auth.HashPassword(*password)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= *n; i++ {
		email := fmt.Sprintf("student%d@toefl.dev", i)
		tag, err := pool.Exec(ctx,
			`INSERT INTO users (email, password_hash, role) VALUES ($1,$2,'student')
			 ON CONFLICT (email) DO NOTHING`, email, hash)
		if err != nil {
			log.Fatalf("%s: %v", email, err)
		}
		log.Printf("ok %s (inserted=%v)", email, tag.RowsAffected() > 0)
	}
}