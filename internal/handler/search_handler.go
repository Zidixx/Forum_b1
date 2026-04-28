package handler

import (
	"encoding/json"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/service"
	"net/http"
)

type SearchHandler struct {
	postService *service.PostService
	tmpl        Renderer
	errHandler  *ErrorHandler
}

func NewSearchHandler(postService *service.PostService, tmpl Renderer, errHandler *ErrorHandler) *SearchHandler {
	return &SearchHandler{postService: postService, tmpl: tmpl, errHandler: errHandler}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	user := middleware.GetUser(r)
	userID := 0
	if user != nil {
		userID = user.ID
	}

	var posts []models.Post
	var err error
	if len(query) >= 2 {
		posts, err = h.postService.Search(query, userID)
	}
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
		"Query": query,
		"Title": "Recherche",
	}

	h.tmpl.ExecuteTemplate(w, "search.html", data)
}

func (h *SearchHandler) APISearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) < 2 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	user := middleware.GetUser(r)
	userID := 0
	if user != nil {
		userID = user.ID
	}

	posts, err := h.postService.Search(query, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal"})
		return
	}

	if len(posts) > 5 {
		posts = posts[:5]
	}

	type searchResult struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Author  string `json:"author"`
		Excerpt string `json:"excerpt"`
	}

	results := make([]searchResult, 0, len(posts))
	for _, p := range posts {
		excerpt := h.postService.Excerpt(p.Content)
		if len(excerpt) > 100 {
			excerpt = excerpt[:100] + "..."
		}
		results = append(results, searchResult{
			ID:      p.ID,
			Title:   p.Title,
			Author:  p.Author,
			Excerpt: excerpt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
