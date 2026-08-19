package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	httppkg "toefl-prep/backend/internal/http"
	"toefl-prep/backend/internal/auth"
)

func CORS(origins []string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, o := range origins {
		allowed[o] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && (len(allowed) == 0 || allowed[origin]) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func AuthRequired(svc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(httppkg.SessionCookie)
		if err != nil {
			httppkg.Fail(c, httppkg.ErrUnauthorized)
			c.Abort()
			return
		}
		user, err := svc.Verify(c.Request.Context(), token)
		if err != nil {
			httppkg.Fail(c, err)
			c.Abort()
			return
		}
		c.Set(httppkg.UserKey, user)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get(httppkg.UserKey)
		if !ok {
			httppkg.Fail(c, httppkg.ErrUnauthorized)
			c.Abort()
			return
		}
		if user.(*auth.User).Role != role {
			httppkg.Fail(c, httppkg.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *auth.User {
	u, _ := c.Get(httppkg.UserKey)
	if u == nil {
		return nil
	}
	return u.(*auth.User)
}

// CSRF applies the double-submit-cookie check to mutating requests.
// The csrf cookie is set at login; the client must echo it in X-CSRF-Token.
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		m := c.Request.Method
		if m == http.MethodGet || m == http.MethodHead || m == http.MethodOptions {
			c.Next()
			return
		}
		cookie, err := c.Cookie(httppkg.CSRFCookie)
		if err != nil || cookie == "" {
			httppkg.Fail(c, httppkg.ErrForbidden)
			c.Abort()
			return
		}
		header := c.GetHeader("X-CSRF-Token")
		if header == "" || !strings.EqualFold(header, cookie) {
			httppkg.Fail(c, httppkg.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

func SetUser(c *gin.Context, user *auth.User) { c.Set(httppkg.UserKey, user) }
func CSRFValue(c *gin.Context) string         { v, _ := c.Get(httppkg.CSRFKey); return v.(string) }
func SetCSRF(c *gin.Context, token string)    { c.Set(httppkg.CSRFKey, token) }