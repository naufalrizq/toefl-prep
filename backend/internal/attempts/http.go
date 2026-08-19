package attempts

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

type startBody struct {
	ExamTemplateID int64  `json:"exam_template_id"`
	Mode           string `json:"mode"`
}

func (h *Handler) Start(c *gin.Context) {
	var body startBody
	if err := httppkg.BindJSON(c, &body); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	result, err := h.svc.Start(c.Request.Context(), user.ID, body.ExamTemplateID, body.Mode)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.Created(c, result)
}

func (h *Handler) List(c *gin.Context) {
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	items, err := h.svc.List(c.Request.Context(), user.ID)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OKList(c, items)
}

func (h *Handler) Questions(c *gin.Context) {
	attemptID, ok := parseID(c)
	if !ok {
		return
	}
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	items, err := h.svc.QuizQuestions(c.Request.Context(), user.ID, attemptID)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OKList(c, items)
}

type answerBody struct {
	ChosenIndex  *int `json:"chosen_index"`
	TimeTakenMs  *int `json:"time_taken_ms"`
}

func (h *Handler) Answer(c *gin.Context) {
	attemptID, ok := parseID(c)
	if !ok {
		return
	}
	itemID, ok := parseItemID(c)
	if !ok {
		return
	}
	var body answerBody
	if err := httppkg.BindJSON(c, &body); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	if err := h.svc.Answer(c.Request.Context(), user.ID, attemptID, itemID, body.ChosenIndex, body.TimeTakenMs); err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.NoContent(c)
}

type flagBody struct {
	Flagged bool `json:"flagged"`
}

func (h *Handler) Flag(c *gin.Context) {
	attemptID, ok := parseID(c)
	if !ok {
		return
	}
	itemID, ok := parseItemID(c)
	if !ok {
		return
	}
	var body flagBody
	if err := httppkg.BindJSON(c, &body); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	if err := h.svc.Flag(c.Request.Context(), user.ID, attemptID, itemID, body.Flagged); err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.NoContent(c)
}

func (h *Handler) Submit(c *gin.Context) {
	attemptID, ok := parseID(c)
	if !ok {
		return
	}
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	result, err := h.svc.Submit(c.Request.Context(), user.ID, attemptID)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, result)
}

func (h *Handler) Review(c *gin.Context) {
	attemptID, ok := parseID(c)
	if !ok {
		return
	}
	user := c.MustGet(httppkg.UserKey).(*auth.User)
	result, err := h.svc.Review(c.Request.Context(), user.ID, attemptID)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	httppkg.OK(c, result)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		httppkg.Fail(c, httppkg.ErrNotFound)
		return 0, false
	}
	return id, true
}

func parseItemID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil || id < 1 {
		httppkg.Fail(c, httppkg.ErrNotFound)
		return 0, false
	}
	return id, true
}