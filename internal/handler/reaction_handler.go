package handler

import (
	"encoding/json"
	"fmt"
	"forum/internal/middleware"
	"forum/internal/repository"
	"forum/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ReactionHandler struct {
	reactionService *service.ReactionService
	repostRepo      *repository.RepostRepository
	errHandler      *ErrorHandler
}

func NewReactionHandler(reactionService *service.ReactionService, repostRepo *repository.RepostRepository, errHandler *ErrorHandler) *ReactionHandler {
	return &ReactionHandler{reactionService: reactionService, repostRepo: repostRepo, errHandler: errHandler}
}

func isAJAX(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
		r.Header.Get("Accept") == "application/json"
}

func (h *ReactionHandler) ReactPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		if isAJAX(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "login_required"})
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/post/react/")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	reactionType := r.FormValue("type")
	if err := h.reactionService.ReactToPost(user.ID, postID, reactionType); err != nil {
		if isAJAX(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "bad_request"})
			return
		}
		h.errHandler.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if isAJAX(r) {
		likes, dislikes, _ := h.reactionService.CountPostReactions(postID)
		userVote, _ := h.reactionService.GetPostReaction(user.ID, postID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"likes":    likes,
			"dislikes": dislikes,
			"userVote": userVote,
		})
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
		if isAJAX(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "login_required"})
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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
		if isAJAX(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "bad_request"})
			return
		}
		h.errHandler.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if isAJAX(r) {
		likes, dislikes, _ := h.reactionService.CountCommentReactions(commentID)
		userVote, _ := h.reactionService.GetCommentReaction(user.ID, commentID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"likes":    likes,
			"dislikes": dislikes,
			"userVote": userVote,
		})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d#comment-%d", postID, commentID), http.StatusSeeOther)
}

func (h *ReactionHandler) Repost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		if isAJAX(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "login_required"})
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/repost/")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	reposted, err := h.repostRepo.Toggle(user.ID, postID)
	if err != nil {
		if isAJAX(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "internal"})
			return
		}
		h.errHandler.InternalError(w, r, err)
		return
	}

	count, _ := h.repostRepo.Count(postID)

	if isAJAX(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"reposted": reposted,
			"count":    count,
		})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
