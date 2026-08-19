package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"toefl-prep/backend/internal/attempts"
	"toefl-prep/backend/internal/auth"
	"toefl-prep/backend/internal/config"
	"toefl-prep/backend/internal/database"
	"toefl-prep/backend/internal/exams"
	"toefl-prep/backend/internal/httpapi"
	"toefl-prep/backend/internal/questions"
	"toefl-prep/backend/internal/reporting"
	"toefl-prep/backend/internal/seed"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := database.Migrate(ctx, cfg.DatabaseURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	pool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	if cfg.SeedUsers {
		if err := seed.EnsureUsers(ctx, pool); err != nil {
			log.Fatalf("seed users: %v", err)
		}
	}

	authSvc := auth.New(pool, cfg.SessionTTL)

	qRepo := questions.NewRepo(pool)
	qSvc := questions.NewService(qRepo)
	qHandler := questions.NewHandler(qSvc)

	seedHandler := seed.NewHandler(qSvc)

	examRepo := exams.NewRepo(pool)
	examSvc := exams.NewService(examRepo, qRepo)
	examHandler := exams.NewHandler(examSvc)

	attemptRepo := attempts.NewRepo(pool)
	attemptSvc := attempts.NewService(attemptRepo, examRepo, qRepo)
	attemptHandler := attempts.NewHandler(attemptSvc)

	reportSvc := reporting.NewService(attemptRepo)
	reportHandler := reporting.NewHandler(reportSvc)

	r := httpapi.New(httpapi.Deps{
		Auth:      auth.NewHandler(authSvc, cfg.LoginRatePerMinute),
		Questions: qHandler,
		Seed:      seedHandler,
		Exams:     examHandler,
		Attempts:  attemptHandler,
		Reporting: reportHandler,
		CORS:      cfg.CORSOrigins,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("toefl-prep backend listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}