package handler

import (
	"forum/internal/middleware"
	"forum/internal/service"
	"html/template"
	"net/http"
)

type ProfileHandler struct {
	postService *service.PostService
	tmpl        *template.Template
	errHandler  *ErrorHandler
}

func NewProfileHandler(postService *service.PostService, tmpl *template.Template, errHandler *ErrorHandler) *ProfileHandler {
	return &ProfileHandler{postService: postService, tmpl: tmpl, errHandler: errHandler}
}

func (h *ProfileHandler) MyPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	posts, err := h.postService.GetByUser(user.ID)
	if err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	var postData []interface{}
	for _, post := range posts {
		postData = append(postData, map[string]interface{}{
			"Post":    post,
			"Excerpt": h.postService.Excerpt(post.Content),
		})
	}

	data := map[string]interface{}{
		"User":  user,
		"Posts": postData,
		"Title": "Mes posts",
	}
	h.tmpl.ExecuteTemplate(w, "my_posts.html", data)
}

func (h *ProfileHandler) LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	posts, err := h.postService.GetLikedByUser(user.ID)
	if err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	var postData []interface{}
	for _, post := range posts {
		postData = append(postData, map[string]interface{}{
			"Post":    post,
			"Excerpt": h.postService.Excerpt(post.Content),
		})
	}

	data := map[string]interface{}{
		"User":  user,
		"Posts": postData,
		"Title": "Posts aimés",
	}
	h.tmpl.ExecuteTemplate(w, "liked_posts.html", data)
}
