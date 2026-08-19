package reporting

import (
	"github.com/gin-gonic/gin"

	"toefl-prep/backend/internal/auth"
	httppkg "toefl-prep/backend/internal/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Dashboard(c *gin.Context) {
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	stats, err := h.svc.Dashboard(c.Request.Context(), user.ID)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, stats)
}