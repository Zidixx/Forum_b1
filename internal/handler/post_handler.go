package handler

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/internal/utils"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	postService   *service.PostService
	uploadService *service.UploadService
	catRepo       *repository.CategoryRepository
	tmpl          *template.Template
	errHandler    *ErrorHandler
}

func NewPostHandler(postService *service.PostService, uploadService *service.UploadService, catRepo *repository.CategoryRepository, tmpl *template.Template, errHandler *ErrorHandler) *PostHandler {
	return &PostHandler{postService: postService, uploadService: uploadService, catRepo: catRepo, tmpl: tmpl, errHandler: errHandler}
}

func (h *PostHandler) ShowCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	categories, _ := h.catRepo.FindAll()
	data := map[string]interface{}{
		"User":       middleware.GetUser(r),
		"Categories": categories,
	}
	h.tmpl.ExecuteTemplate(w, "create_post.html", data)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	r.ParseMultipartForm(utils.MaxUploadSize)

	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryStrs := r.Form["categories"]

	var categoryIDs []int
	for _, s := range categoryStrs {
		id, err := strconv.Atoi(s)
		if err == nil {
			categoryIDs = append(categoryIDs, id)
		}
	}

	errs := utils.ValidatePost(title, content, categoryIDs)

	// Handle image upload
	var imagePath string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		filename, uploadErr := h.uploadService.SaveImage(file, header)
		if uploadErr != nil {
			errs["image"] = uploadErr.Error()
		} else {
			imagePath = filename
		}
	}

	if errs.HasErrors() {
		categories, _ := h.catRepo.FindAll()
		data := map[string]interface{}{
			"User":        user,
			"Errors":      errs,
			"Title":       title,
			"Content":     content,
			"Categories":  categories,
			"SelectedCats": categoryIDs,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "create_post.html", data)
		return
	}

	post := &models.Post{
		UserID:    user.ID,
		Title:     strings.TrimSpace(title),
		Content:   strings.TrimSpace(content),
		ImagePath: imagePath,
	}

	if err := h.postService.Create(post, categoryIDs); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusSeeOther)
}

func (h *PostHandler) Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	id, err := h.extractID(r.URL.Path, "/post/")
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	user := middleware.GetUser(r)
	userID := 0
	if user != nil {
		userID = user.ID
	}

	post, err := h.postService.GetByID(id, userID)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"User": user,
		"Post": post,
	}

	h.tmpl.ExecuteTemplate(w, "post_detail.html", data)
}

func (h *PostHandler) ShowEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := h.extractID(r.URL.Path, "/post/edit/")
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	post, err := h.postService.GetByID(id, user.ID)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	if post.UserID != user.ID {
		h.errHandler.Forbidden(w, r)
		return
	}

	catIDs, _ := h.postService.GetCategoryIDs(post.ID)
	categories, _ := h.catRepo.FindAll()

	data := map[string]interface{}{
		"User":        user,
		"Post":        post,
		"Categories":  categories,
		"SelectedCats": catIDs,
	}

	h.tmpl.ExecuteTemplate(w, "edit_post.html", data)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := h.extractID(r.URL.Path, "/post/edit/")
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	post, err := h.postService.GetByID(id, user.ID)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	if post.UserID != user.ID {
		h.errHandler.Forbidden(w, r)
		return
	}

	r.ParseMultipartForm(utils.MaxUploadSize)

	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryStrs := r.Form["categories"]

	var categoryIDs []int
	for _, s := range categoryStrs {
		cid, convErr := strconv.Atoi(s)
		if convErr == nil {
			categoryIDs = append(categoryIDs, cid)
		}
	}

	errs := utils.ValidatePost(title, content, categoryIDs)

	// Handle new image
	file, header, fileErr := r.FormFile("image")
	if fileErr == nil {
		defer file.Close()
		filename, uploadErr := h.uploadService.SaveImage(file, header)
		if uploadErr != nil {
			errs["image"] = uploadErr.Error()
		} else {
			post.ImagePath = filename
		}
	}

	if errs.HasErrors() {
		categories, _ := h.catRepo.FindAll()
		data := map[string]interface{}{
			"User":        user,
			"Post":        post,
			"Errors":      errs,
			"Categories":  categories,
			"SelectedCats": categoryIDs,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "edit_post.html", data)
		return
	}

	post.Title = strings.TrimSpace(title)
	post.Content = strings.TrimSpace(content)

	if err := h.postService.Update(post, categoryIDs); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusSeeOther)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errHandler.MethodNotAllowed(w, r)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := h.extractID(r.URL.Path, "/post/delete/")
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	post, err := h.postService.GetByID(id, user.ID)
	if err != nil {
		h.errHandler.NotFound(w, r)
		return
	}

	if post.UserID != user.ID {
		h.errHandler.Forbidden(w, r)
		return
	}

	if err := h.postService.Delete(id); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PostHandler) extractID(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	return strconv.Atoi(idStr)
}
