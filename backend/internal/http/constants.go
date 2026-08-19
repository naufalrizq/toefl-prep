package http

const (
	// SessionCookie is the httpOnly session cookie (double-submit CSRF uses
	// the sibling csrf cookie). Names live here so both auth and middleware
	// can reference them without an import cycle.
	SessionCookie = "toefl_session"
	CSRFCookie    = "toefl_csrf"
	UserKey       = "auth_user"
	CSRFKey       = "csrf_token"
)