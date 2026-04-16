package middleware

import (
	"forum/internal/service"
	"net/http"
)

// RedirectIfAuth redirects authenticated users away from login/register pages
func RedirectIfAuth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err == nil && cookie.Value != "" {
				user, err := authService.GetUserFromSession(cookie.Value)
				if err == nil && user != nil {
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
