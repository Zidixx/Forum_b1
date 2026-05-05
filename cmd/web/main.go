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

// League represents a football league for template rendering
type League struct {
	Slug string
	Name string
	Flag string
}

var leagues = []League{
	{"ligue1", "Ligue 1", "\U0001F1EB\U0001F1F7"},
	{"premier-league", "Premier League", "\U0001F3F4\U000E0067\U000E0062\U000E0065\U000E006E\U000E0067\U000E007F"},
	{"la-liga", "La Liga", "\U0001F1EA\U0001F1F8"},
	{"bundesliga", "Bundesliga", "\U0001F1E9\U0001F1EA"},
	{"serie-a", "Serie A", "\U0001F1EE\U0001F1F9"},
	{"champions-league", "Champions League", "\U0001F3C6"},
	{"europa-league", "Europa League", "\U0001F3C6"},
}

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
		"lower":   strings.ToLower,
		"add": func(a, b int) int {
			return a + b
		},
		"contains": func(slice []int, val int) bool {
			for _, v := range slice {
				if v == val {
					return true
				}
			}
			return false
		},
		"leagues": func() []League {
			return leagues
		},
		"leagueLabel": func(slug string) string {
			for _, l := range leagues {
				if l.Slug == slug {
					return l.Name
				}
			}
			return slug
		},
		"leagueFlag": func(slug string) string {
			for _, l := range leagues {
				if l.Slug == slug {
					return l.Flag
				}
			}
			return ""
		},
		"leagueColor": func(slug string) string {
			colors := map[string]string{
				"ligue1":           "#091c3e",
				"premier-league":   "#3d195b",
				"la-liga":          "#ee8707",
				"bundesliga":       "#d20515",
				"serie-a":          "#024494",
				"champions-league": "#1a3a5c",
				"europa-league":    "#f47920",
			}
			if c, ok := colors[slug]; ok {
				return c
			}
			return "#333333"
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
	homeHandler := handler.NewHomeHandler(postService, catRepo, userRepo, commentRepo, tmpl, errHandler)
	postHandler := handler.NewPostHandler(postService, uploadService, commentService, catRepo, tmpl, errHandler)
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

	fmt.Printf("\n  LE VESTIAIRE - Forum Football\n")
	fmt.Printf("  Démarré en HTTPS sur https://localhost:%s\n", port)
	fmt.Printf("  Rate Limiting Anti-DDoS actif\n\n")

	err = http.ListenAndServeTLS(":"+port, "./tls/server.crt", "./tls/server.key", protectedMux)
	if err != nil {
		log.Fatal(err)
	}
}
