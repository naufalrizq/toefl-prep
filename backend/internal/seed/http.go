package seed

import (
	"github.com/gin-gonic/gin"

	httppkg "toefl-prep/backend/internal/http"
	"toefl-prep/backend/internal/questions"
)

type Handler struct {
	svc *questions.Service
}

func NewHandler(svc *questions.Service) *Handler { return &Handler{svc: svc} }

// Seed loads the embedded question bank and upserts it (admin only, FR-2.8).
func (h *Handler) Seed(c *gin.Context) {
	qs, err := Load()
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	n, err := h.svc.Seed(c.Request.Context(), qs)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, gin.H{"seeded": n})
}