package handler

import (
	"forum/internal/middleware"
	"forum/internal/repository"
	"forum/internal/service"
	"net/http"
	"strconv"
)

type HomeHandler struct {
	postService *service.PostService
	catRepo     *repository.CategoryRepository
	userRepo    *repository.UserRepository
	commentRepo *repository.CommentRepository
	tmpl        Renderer
	errHandler  *ErrorHandler
}

func NewHomeHandler(postService *service.PostService, catRepo *repository.CategoryRepository, userRepo *repository.UserRepository, commentRepo *repository.CommentRepository, tmpl Renderer, errHandler *ErrorHandler) *HomeHandler {
	return &HomeHandler{postService: postService, catRepo: catRepo, userRepo: userRepo, commentRepo: commentRepo, tmpl: tmpl, errHandler: errHandler}
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
	leagueFilter := r.URL.Query().Get("league")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "new"
	}

	var rawPosts []interface{}
	var err error

	if leagueFilter != "" {
		p, fetchErr := h.postService.GetByLeague(leagueFilter, userID)
		if fetchErr != nil {
			h.errHandler.InternalError(w, r, fetchErr)
			return
		}
		for _, post := range p {
			rawPosts = append(rawPosts, map[string]interface{}{
				"Post":    post,
				"Excerpt": h.postService.Excerpt(post.Content),
			})
		}
	} else if categoryFilter != "" {
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
			rawPosts = append(rawPosts, map[string]interface{}{
				"Post":    post,
				"Excerpt": h.postService.Excerpt(post.Content),
			})
		}
	} else {
		p, fetchErr := h.postService.GetAllSorted(sort, userID)
		if fetchErr != nil {
			h.errHandler.InternalError(w, r, fetchErr)
			return
		}
		for _, post := range p {
			rawPosts = append(rawPosts, map[string]interface{}{
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

	trending, _ := h.postService.GetTrending(5, userID)

	totalPosts, _ := h.postService.CountAll()
	totalUsers, _ := h.userRepo.CountAll()
	totalComments, _ := h.commentRepo.CountAll()

	data := map[string]interface{}{
		"User":           user,
		"Posts":          rawPosts,
		"Categories":     categories,
		"CategoryFilter": categoryFilter,
		"LeagueFilter":   leagueFilter,
		"Sort":           sort,
		"Trending":       trending,
		"TotalPosts":     totalPosts,
		"TotalUsers":     totalUsers,
		"TotalComments":  totalComments,
	}

	h.tmpl.ExecuteTemplate(w, "home.html", data)
}
