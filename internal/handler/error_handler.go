package handler

import (
	"forum/internal/middleware"
	"log"
	"net/http"
)

type ErrorHandler struct {
	tmpl Renderer
}

func NewErrorHandler(tmpl Renderer) *ErrorHandler {
	return &ErrorHandler{tmpl: tmpl}
}

var errorMessages = map[int]string{
	404: "HORS-JEU ! Cette page est sortie du terrain.",
	403: "CARTON ROUGE ! Tu n'as pas ton pass VIP pour cette zone.",
	405: "MAUVAISE PASSE ! Cette action n'est pas permise ici.",
	500: "BLESSURE DU SERVEUR ! Notre gardien a laissé passer une erreur.",
	429: "COUP DE SIFFLET ! Trop de requêtes, l'arbitre siffle la pause.",
}

func (h *ErrorHandler) Error(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	footballMsg := errorMessages[status]
	if footballMsg == "" {
		footballMsg = message
	}
	data := map[string]interface{}{
		"User":        middleware.GetUser(r),
		"StatusCode":  status,
		"Message":     footballMsg,
		"SubMessage":  message,
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
	h.Error(w, r, http.StatusMethodNotAllowed, "Méthode non autoris��e")
}

func (h *ErrorHandler) InternalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error: %v", err)
	h.Error(w, r, http.StatusInternalServerError, "Erreur interne du serveur")
}

func (h *ErrorHandler) Forbidden(w http.ResponseWriter, r *http.Request) {
	h.Error(w, r, http.StatusForbidden, "Accès interdit")
}
