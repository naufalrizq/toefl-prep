package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	httppkg "toefl-prep/backend/internal/http"
)

// loginLimiter is a tiny per-IP sliding window used to slow brute force on
// the login endpoint. In-memory by design: single instance, fine for a
// personal app; swap for a Redis-backed limiter if the app grows.
type loginLimiter struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	seen    map[string][]time.Time
	cleanup time.Time
}

func newLoginLimiter(max int, window time.Duration) *loginLimiter {
	return &loginLimiter{
		window:  window,
		max:     max,
		seen:    make(map[string][]time.Time),
		cleanup: time.Now(),
	}
}

func (l *loginLimiter) allow(key string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if now.After(l.cleanup) {
		for ip, hits := range l.seen {
			var kept []time.Time
			for _, t := range hits {
				if now.Sub(t) <= l.window {
					kept = append(kept, t)
				}
			}
			if len(kept) == 0 {
				delete(l.seen, ip)
			} else {
				l.seen[ip] = kept
			}
		}
		l.cleanup = now.Add(l.window)
	}
	hits := l.seen[key]
	if len(hits) >= l.max {
		return false
	}
	l.seen[key] = append(hits, now)
	return true
}

type Handler struct {
	svc     *Service
	limiter *loginLimiter
}

func NewHandler(svc *Service, loginPerMinute int) *Handler {
	return &Handler{svc: svc, limiter: newLoginLimiter(loginPerMinute, time.Minute)}
}

func (h *Handler) Service() *Service { return h.svc }

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	ip := c.ClientIP()
	if !h.limiter.allow(ip) {
		httppkg.Fail(c, httppkg.NewError(429, "rate_limited", "too many login attempts, try again later"))
		return
	}
	var body loginBody
	if err := httppkg.BindJSON(c, &body); err != nil {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", err.Error()))
		return
	}
	if body.Email == "" || body.Password == "" {
		httppkg.Fail(c, httppkg.NewError(422, "validation_failed", "email and password are required"))
		return
	}

	token, err := h.svc.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}
	user, err := h.svc.Verify(c.Request.Context(), token)
	if err != nil {
		httppkg.Fail(c, err)
		return
	}

	csrf, err := csrfToken()
	if err != nil {
		httppkg.Fail(c, err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     httppkg.SessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.svc.ttl.Seconds()),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     httppkg.CSRFCookie,
		Value:    csrf,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.svc.ttl.Seconds()),
	})
	httppkg.OK(c, user)
}

func (h *Handler) Me(c *gin.Context) {
	user := c.MustGet(httppkg.UserKey).(*User)
	httppkg.OK(c, user)
}

func (h *Handler) Logout(c *gin.Context) {
	token, err := c.Cookie(httppkg.SessionCookie)
	if err == nil && token != "" {
		_ = h.svc.Logout(c.Request.Context(), token)
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     httppkg.SessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     httppkg.CSRFCookie,
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	httppkg.NoContent(c)
}

func csrfToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}