package exams

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"toefl-prep/backend/internal/auth"
	httppkg "toefl-prep/backend/internal/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// List serves both roles on GET /exams: admins see every template, students
// only published/active ones. One route keeps the router free of duplicates.
func (h *Handler) List(c *gin.Context) {
	if user := c.MustGet(httppkg.UserKey).(*auth.User); user.Role == "admin" {
		h.ListAdmin(c)
		return
	}
	h.ListPublished(c)
}

type publishBody struct {
	Published bool `json:"published"`
}

func (h *Handler) ListAdmin(c *gin.Context) {
	items, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OKList(c, items)
}

func (h *Handler) ListPublished(c *gin.Context) {
	items, err := h.svc.ListPublished(c.Request.Context())
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OKList(c, items)
}

func (h *Handler) Create(c *gin.Context) {
	var e ExamTemplate
	if err := httppkg.BindJSON(c, &e); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	if err := h.svc.Validate(&e); err != nil {
		httppkg.Fail(c, err)
		return
	}
	id, err := h.svc.Create(c.Request.Context(), &e)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	created, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.Created(c, created)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var e ExamTemplate
	if err := httppkg.BindJSON(c, &e); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	e.ID = id
	if err := h.svc.Validate(&e); err != nil {
		httppkg.Fail(c, err)
		return
	}
	if err := h.svc.Update(c.Request.Context(), &e); err != nil {
		httppkg.Fail(c, err)
		return
	}
	updated, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.SoftDelete(c.Request.Context(), id); err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.NoContent(c)
}

func (h *Handler) Publish(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var body publishBody
	if err := httppkg.BindJSON(c, &body); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	if err := h.svc.Publish(c.Request.Context(), id, body.Published); err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.NoContent(c)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		httppkg.Fail(c, httppkg.ErrNotFound)
		return 0, false
	}
	return id, true
}