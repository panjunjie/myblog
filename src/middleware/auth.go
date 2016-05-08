package middleware

import (
	"myblog/src/controller"
	"net/http"
	"strings"
)

type CheckLogin struct {
	LoginUrl     string
	RememberNext bool
}

// NewLogger returns a new Logger instance
func NewCheckLogin() *CheckLogin {
	return &CheckLogin{"/account/login", true}
}

func (t *CheckLogin) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if controller.GetSession(req, controller.SESSION_WEB) == nil {
		var login_url = t.LoginUrl
		next := req.URL.Path
		next_lower := strings.ToLower(next)
		if t.RememberNext && (!strings.Contains(next_lower, "login") || !strings.Contains(next_lower, "logout")) {
			login_url = login_url + "?next=" + next
		}

		http.Redirect(rw, req, login_url, http.StatusFound)
		return
	} else {
		next(rw, req)
	}
}
