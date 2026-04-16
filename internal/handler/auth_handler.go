package handler

import (
	"forum/internal/middleware"
	"forum/internal/service"
	"forum/internal/utils"
	"html/template"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService *service.AuthService
	tmpl        *template.Template
	errHandler  *ErrorHandler
}

func NewAuthHandler(authService *service.AuthService, tmpl *template.Template, errHandler *ErrorHandler) *AuthHandler {
	return &AuthHandler{authService: authService, tmpl: tmpl, errHandler: errHandler}
}

func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}
	data := map[string]interface{}{
		"User": middleware.GetUser(r),
	}
	h.tmpl.ExecuteTemplate(w, "register.html", data)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	errs := utils.ValidateRegister(email, username, password, confirmPassword)
	if errs.HasErrors() {
		data := map[string]interface{}{
			"Errors":   errs,
			"Email":    email,
			"Username": username,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "register.html", data)
		return
	}

	_, serviceErrs := h.authService.Register(email, username, password)
	if serviceErrs != nil && serviceErrs.HasErrors() {
		data := map[string]interface{}{
			"Errors":   serviceErrs,
			"Email":    email,
			"Username": username,
		}
		w.WriteHeader(http.StatusConflict)
		h.tmpl.ExecuteTemplate(w, "register.html", data)
		return
	}

	http.Redirect(w, r, "/login?registered=1", http.StatusSeeOther)
}

func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}
	data := map[string]interface{}{
		"User":       middleware.GetUser(r),
		"Registered": r.URL.Query().Get("registered") == "1",
	}
	h.tmpl.ExecuteTemplate(w, "login.html", data)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	errs := utils.ValidateLogin(identifier, password)
	if errs.HasErrors() {
		data := map[string]interface{}{
			"Errors":     errs,
			"Identifier": identifier,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "login.html", data)
		return
	}

	_, session, err := h.authService.Login(identifier, password)
	if err != nil {
		data := map[string]interface{}{
			"Error":      err.Error(),
			"Identifier": identifier,
		}
		w.WriteHeader(http.StatusUnauthorized)
		h.tmpl.ExecuteTemplate(w, "login.html", data)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.authService.Logout(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
