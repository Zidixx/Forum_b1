package handler

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/service"
	"forum/internal/utils"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	commentService *service.CommentService
	tmpl           *template.Template
	errHandler     *ErrorHandler
}

func NewCommentHandler(commentService *service.CommentService, tmpl *template.Template, errHandler *ErrorHandler) *CommentHandler {
	return &CommentHandler{commentService: commentService, tmpl: tmpl, errHandler: errHandler}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.errHandler.Error(w, r, http.StatusBadRequest, "Post invalide")
		return
	}

	content := r.FormValue("content")
	errs := utils.ValidateComment(content)
	if errs.HasErrors() {
		http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
		return
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  user.ID,
		Content: strings.TrimSpace(content),
	}

	if err := h.commentService.Create(comment); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d#comment-%d", postID, comment.ID), http.StatusSeeOther)
}

func (h *CommentHandler) ShowEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/comment/edit/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	comment, err := h.commentService.GetByID(id)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	if comment.UserID != user.ID {
		h.errHandler.Forbidden(w, r)
		return
	}

	data := map[string]interface{}{
		"User":    user,
		"Comment": comment,
	}
	h.tmpl.ExecuteTemplate(w, "edit_comment.html", data)
}

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/comment/edit/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	comment, err := h.commentService.GetByID(id)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	if comment.UserID != user.ID {
		h.errHandler.Forbidden(w, r)
		return
	}

	content := r.FormValue("content")
	errs := utils.ValidateComment(content)
	if errs.HasErrors() {
		data := map[string]interface{}{
			"User":    user,
			"Comment": comment,
			"Errors":  errs,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "edit_comment.html", data)
		return
	}

	comment.Content = strings.TrimSpace(content)
	if err := h.commentService.Update(comment); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d#comment-%d", comment.PostID, comment.ID), http.StatusSeeOther)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/comment/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	comment, err := h.commentService.GetByID(id)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	if comment.UserID != user.ID {
		h.errHandler.Forbidden(w, r)
		return
	}

	postID := comment.PostID
	if err := h.commentService.Delete(id); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
