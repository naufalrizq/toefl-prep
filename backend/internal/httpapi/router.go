// Package httpapi assembles the HTTP surface. It lives outside internal/http
// (which owns the shared envelope/errors and is imported by every service)
// so the router can depend on auth + middleware + handlers without a cycle.
package httpapi

import (
	"github.com/gin-gonic/gin"

	"toefl-prep/backend/internal/attempts"
	"toefl-prep/backend/internal/auth"
	"toefl-prep/backend/internal/exams"
	"toefl-prep/backend/internal/http/middleware"
	httppkg "toefl-prep/backend/internal/http"
	"toefl-prep/backend/internal/questions"
	"toefl-prep/backend/internal/reporting"
	"toefl-prep/backend/internal/seed"
)

// Deps wires every handler for the API surface described in SRS §8.
type Deps struct {
	Auth      *auth.Handler
	Questions *questions.Handler
	Seed      *seed.Handler
	Exams     *exams.Handler
	Attempts  *attempts.Handler
	Reporting *reporting.Handler
	CORS      []string
}

// New builds the router with all middleware and routes.
func New(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(d.CORS))

	v1 := r.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		httppkg.OK(c, gin.H{"status": "ok"})
	})

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", d.Auth.Login)
		authGroup.GET("/me", middleware.AuthRequired(d.Auth.Service()), d.Auth.Me)
		authGroup.POST("/logout", middleware.AuthRequired(d.Auth.Service()), middleware.CSRF(), d.Auth.Logout)
	}

	authed := v1.Group("", middleware.AuthRequired(d.Auth.Service()), middleware.CSRF())
	{
		authed.GET("/exams", d.Exams.List)
		authed.GET("/attempts", d.Attempts.List)
		authed.POST("/attempts", d.Attempts.Start)
		authed.GET("/attempts/:id/questions", d.Attempts.Questions)
		authed.PUT("/attempts/:id/answers/:item_id", d.Attempts.Answer)
		authed.PUT("/attempts/:id/flag/:item_id", d.Attempts.Flag)
		authed.POST("/attempts/:id/submit", d.Attempts.Submit)
		authed.GET("/attempts/:id/review", d.Attempts.Review)
		authed.GET("/dashboard/stats", d.Reporting.Dashboard)
	}

	admin := authed.Group("", middleware.RequireRole("admin"))
	{
		admin.GET("/questions", d.Questions.List)
		admin.POST("/questions", d.Questions.Create)
		admin.GET("/questions/:id", d.Questions.Get)
		admin.PUT("/questions/:id", d.Questions.Update)
		admin.DELETE("/questions/:id", d.Questions.Delete)
		admin.POST("/questions/import", d.Questions.Import)
		admin.POST("/seed", d.Seed.Seed)
		admin.POST("/exams", d.Exams.Create)
		admin.PUT("/exams/:id", d.Exams.Update)
		admin.DELETE("/exams/:id", d.Exams.Delete)
		admin.POST("/exams/:id/publish", d.Exams.Publish)
	}

	return r
}