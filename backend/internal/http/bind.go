package http

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
)

// BindJSON decodes a request body, rejecting unknown fields and trailing
// content so typos in JSON never silently pass.
func BindJSON(c *gin.Context, dst any) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := dec.Decode(&extra); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}