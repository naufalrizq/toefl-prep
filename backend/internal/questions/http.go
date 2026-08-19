package questions

import (
	"strconv"

	"github.com/gin-gonic/gin"

	httppkg "toefl-prep/backend/internal/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		httppkg.Fail(c, httppkg.ErrNotFound)
		return 0, false
	}
	return id, true
}

func bindJSON(c *gin.Context, dst any) bool {
	if err := httppkg.BindJSON(c, dst); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return false
	}
	return true
}

func (h *Handler) List(c *gin.Context) {
	f := Filter{
		Section:    c.Query("section"),
		Type:       c.Query("type"),
		Difficulty: c.Query("difficulty"),
		Search:     c.Query("search"),
	}
	f.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	f.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := h.svc.Repo().List(c.Request.Context(), f)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, gin.H{
		"items": items,
		"page":  f.Page,
		"limit": f.Limit,
		"total": total,
	})
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	q, err := h.svc.Repo().GetByID(c.Request.Context(), id)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, q)
}

func (h *Handler) Create(c *gin.Context) {
	var q Question
	if !bindJSON(c, &q) {
		return
	}
	if err := Validate(&q); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	id, err := h.svc.Repo().Create(c.Request.Context(), &q)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	created, err := h.svc.Repo().GetByID(c.Request.Context(), id)
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
	var q Question
	if !bindJSON(c, &q) {
		return
	}
	q.ID = id
	if err := Validate(&q); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	if err := h.svc.Repo().Update(c.Request.Context(), &q); err != nil {
		httppkg.Fail(c, err)
		return
	}
	updated, err := h.svc.Repo().GetByID(c.Request.Context(), id)
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
	if err := h.svc.Repo().SoftDelete(c.Request.Context(), id); err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.NoContent(c)
}

func (h *Handler) Import(c *gin.Context) {
	results, err := h.svc.Import(c.Request.Context(), c.Request.Body)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, results)
}