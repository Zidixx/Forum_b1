package main

import (
	"fmt"
	"forum/internal/db"
	"forum/internal/handler"
	"forum/internal/middleware"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/internal/utils"
	htmltemplate "html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	baseDir := "."
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	}

	dbPath := filepath.Join(baseDir, "data", "forum.db")
	schemaPath := filepath.Join(baseDir, "sql", "schema.sql")
	seedPath := filepath.Join(baseDir, "sql", "seed.sql")
	templatesDir := filepath.Join(baseDir, "templates")
	staticDir := filepath.Join(baseDir, "static")
	uploadDir := filepath.Join(staticDir, "uploads")

	// Database
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("database open: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database, schemaPath); err != nil {
		log.Fatalf("migration: %v", err)
	}
	if err := db.Seed(database, seedPath); err != nil {
		log.Fatalf("seed: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(database)
	sessionRepo := repository.NewSessionRepository(database)
	postRepo := repository.NewPostRepository(database)
	commentRepo := repository.NewCommentRepository(database)
	catRepo := repository.NewCategoryRepository(database)
	reactionRepo := repository.NewReactionRepository(database)
	repostRepo := repository.NewRepostRepository(database)

	// Services
	authService := service.NewAuthService(userRepo, sessionRepo)
	postService := service.NewPostService(postRepo, catRepo, reactionRepo, repostRepo)
	commentService := service.NewCommentService(commentRepo, reactionRepo)
	reactionService := service.NewReactionService(reactionRepo)
	uploadService := service.NewUploadService(uploadDir)

	// Template functions
	funcMap := htmltemplate.FuncMap{
		"timeAgo": utils.TimeAgo,
		"upper":   strings.ToUpper,
		"add": func(a, b int) int {
			return a + b
		},
	}

	// Templates
	tmpl, err := handler.NewTemplateMap(templatesDir, funcMap)
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	// Handlers
	errHandler := handler.NewErrorHandler(tmpl)
	authHandler := handler.NewAuthHandler(authService, tmpl, errHandler)
	homeHandler := handler.NewHomeHandler(postService, catRepo, tmpl, errHandler)
	postHandler := handler.NewPostHandler(postService, uploadService, catRepo, tmpl, errHandler)
	commentHandler := handler.NewCommentHandler(commentService, tmpl, errHandler)
	reactionHandler := handler.NewReactionHandler(reactionService, repostRepo, errHandler)
	profileHandler := handler.NewProfileHandler(postService, tmpl, errHandler)
	searchHandler := handler.NewSearchHandler(postService, tmpl, errHandler)

	// Middleware
	userCtx := middleware.UserContext(authService)
	requireAuth := middleware.RequireAuth(authService)
	guestOnly := middleware.RedirectIfAuth(authService)

	// Router
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public routes (with user context)
	mux.Handle("/", userCtx(http.HandlerFunc(homeHandler.Home)))
	mux.Handle("/post/", userCtx(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/post/create" && r.Method == http.MethodGet:
			requireAuth(http.HandlerFunc(postHandler.ShowCreate)).ServeHTTP(w, r)
		case path == "/post/create" && r.Method == http.MethodPost:
			requireAuth(http.HandlerFunc(postHandler.Create)).ServeHTTP(w, r)
		case len(path) > len("/post/edit/") && path[:len("/post/edit/")] == "/post/edit/" && r.Method == http.MethodGet:
			requireAuth(http.HandlerFunc(postHandler.ShowEdit)).ServeHTTP(w, r)
		case len(path) > len("/post/edit/") && path[:len("/post/edit/")] == "/post/edit/" && r.Method == http.MethodPost:
			requireAuth(http.HandlerFunc(postHandler.Update)).ServeHTTP(w, r)
		case len(path) > len("/post/delete/") && path[:len("/post/delete/")] == "/post/delete/":
			requireAuth(http.HandlerFunc(postHandler.Delete)).ServeHTTP(w, r)
		case len(path) > len("/post/react/") && path[:len("/post/react/")] == "/post/react/":
			requireAuth(http.HandlerFunc(reactionHandler.ReactPost)).ServeHTTP(w, r)
		default:
			postHandler.Show(w, r)
		}
	})))

	// Comment routes
	mux.Handle("/comment/", userCtx(requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/comment/create":
			commentHandler.Create(w, r)
		case len(path) > len("/comment/edit/") && path[:len("/comment/edit/")] == "/comment/edit/" && r.Method == http.MethodGet:
			commentHandler.ShowEdit(w, r)
		case len(path) > len("/comment/edit/") && path[:len("/comment/edit/")] == "/comment/edit/" && r.Method == http.MethodPost:
			commentHandler.Update(w, r)
		case len(path) > len("/comment/delete/") && path[:len("/comment/delete/")] == "/comment/delete/":
			commentHandler.Delete(w, r)
		case len(path) > len("/comment/react/") && path[:len("/comment/react/")] == "/comment/react/":
			reactionHandler.ReactComment(w, r)
		default:
			errHandler.NotFound(w, r)
		}
	}))))

	// Repost route
	mux.Handle("/repost/", userCtx(requireAuth(http.HandlerFunc(reactionHandler.Repost))))

	// Search routes
	mux.Handle("/search", userCtx(http.HandlerFunc(searchHandler.Search)))
	mux.Handle("/api/search", userCtx(http.HandlerFunc(searchHandler.APISearch)))

	// Auth routes
	mux.Handle("/register", userCtx(guestOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authHandler.ShowRegister(w, r)
		case http.MethodPost:
			authHandler.Register(w, r)
		default:
			errHandler.MethodNotAllowed(w, r)
		}
	}))))

	mux.Handle("/login", userCtx(guestOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authHandler.ShowLogin(w, r)
		case http.MethodPost:
			authHandler.Login(w, r)
		default:
			errHandler.MethodNotAllowed(w, r)
		}
	}))))

	mux.Handle("/logout", userCtx(http.HandlerFunc(authHandler.Logout)))

	mux.Handle("/my-posts", userCtx(requireAuth(http.HandlerFunc(profileHandler.MyPosts))))
	mux.Handle("/liked-posts", userCtx(requireAuth(http.HandlerFunc(profileHandler.LikedPosts))))

	port := "8443"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	protectedMux := middleware.RateLimiter(mux)

	fmt.Printf("Forum démarré en HTTPS sur https://localhost:%s (Avec Rate Limiting Anti-DDoS)\n", port)

	err = http.ListenAndServeTLS(":"+port, "./tls/server.crt", "./tls/server.key", protectedMux)
	if err != nil {
		log.Fatal(err)
	}
}
