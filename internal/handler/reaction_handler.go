package handler

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ReactionHandler struct {
	reactionService *service.ReactionService
	errHandler      *ErrorHandler
}

func NewReactionHandler(reactionService *service.ReactionService, errHandler *ErrorHandler) *ReactionHandler {
	return &ReactionHandler{reactionService: reactionService, errHandler: errHandler}
}

func (h *ReactionHandler) ReactPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// /post/react/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/post/react/")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	reactionType := r.FormValue("type")
	if err := h.reactionService.ReactToPost(user.ID, postID, reactionType); err != nil {
		h.errHandler.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

func (h *ReactionHandler) ReactComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// /comment/react/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/comment/react/")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	reactionType := r.FormValue("type")
	postIDStr := r.FormValue("post_id")
	postID, _ := strconv.Atoi(postIDStr)

	if err := h.reactionService.ReactToComment(user.ID, commentID, reactionType); err != nil {
		h.errHandler.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d#comment-%d", postID, commentID), http.StatusSeeOther)
}
