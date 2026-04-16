package handler

import (
	"forum/internal/middleware"
	"html/template"
	"log"
	"net/http"
)

type ErrorHandler struct {
	tmpl *template.Template
}

func NewErrorHandler(tmpl *template.Template) *ErrorHandler {
	return &ErrorHandler{tmpl: tmpl}
}

func (h *ErrorHandler) Error(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	data := map[string]interface{}{
		"User":       middleware.GetUser(r),
		"StatusCode": status,
		"Message":    message,
	}
	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, message, status)
	}
}

func (h *ErrorHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.Error(w, r, http.StatusNotFound, "Page introuvable")
}

func (h *ErrorHandler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	h.Error(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
}

func (h *ErrorHandler) InternalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error: %v", err)
	h.Error(w, r, http.StatusInternalServerError, "Erreur interne du serveur")
}

func (h *ErrorHandler) Forbidden(w http.ResponseWriter, r *http.Request) {
	h.Error(w, r, http.StatusForbidden, "Accès interdit")
}
