package handler

import (
	"forum/internal/middleware"
	"forum/internal/repository"
	"forum/internal/service"
	"html/template"
	"net/http"
	"strconv"
)

type HomeHandler struct {
	postService *service.PostService
	catRepo     *repository.CategoryRepository
	tmpl        *template.Template
	errHandler  *ErrorHandler
}

func NewHomeHandler(postService *service.PostService, catRepo *repository.CategoryRepository, tmpl *template.Template, errHandler *ErrorHandler) *HomeHandler {
	return &HomeHandler{postService: postService, catRepo: catRepo, tmpl: tmpl, errHandler: errHandler}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errHandler.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	userID := 0
	if user != nil {
		userID = user.ID
	}

	categoryFilter := r.URL.Query().Get("category")

	var posts []interface{}
	var err error
	var rawPosts []interface{}
	_ = rawPosts

	if categoryFilter != "" {
		catID, convErr := strconv.Atoi(categoryFilter)
		if convErr != nil {
			h.errHandler.Error(w, r, http.StatusBadRequest, "Catégorie invalide")
			return
		}
		p, fetchErr := h.postService.GetByCategory(catID, userID)
		if fetchErr != nil {
			h.errHandler.InternalError(w, r, fetchErr)
			return
		}
		for _, post := range p {
			posts = append(posts, map[string]interface{}{
				"Post":    post,
				"Excerpt": h.postService.Excerpt(post.Content),
			})
		}
	} else {
		p, fetchErr := h.postService.GetAll(userID)
		if fetchErr != nil {
			h.errHandler.InternalError(w, r, fetchErr)
			return
		}
		for _, post := range p {
			posts = append(posts, map[string]interface{}{
				"Post":    post,
				"Excerpt": h.postService.Excerpt(post.Content),
			})
		}
	}

	categories, err := h.catRepo.FindAll()
	if err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	data := map[string]interface{}{
		"User":           user,
		"Posts":          posts,
		"Categories":     categories,
		"CategoryFilter": categoryFilter,
	}

	h.tmpl.ExecuteTemplate(w, "home.html", data)
}
