package http

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Data any `json:"data,omitempty"`
}

type errEnvelope struct {
	Error errBody `json:"error"`
}

type errBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, envelope{Data: data})
}

// OKList marshals a slice, always emitting `[]` instead of `null` for
// empty/nil results so consumers can rely on array envelopes.
func OKList(c *gin.Context, slice any) {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Slice && v.IsNil() {
		slice = reflect.MakeSlice(v.Type(), 0, 0).Interface()
	}
	c.JSON(http.StatusOK, envelope{Data: slice})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, envelope{Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, err error) {
	status, code, message := Map(err)
	c.JSON(status, errEnvelope{Error: errBody{Code: code, Message: message}})
}