package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/internal/utils"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	authService *service.AuthService
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	errHandler  *ErrorHandler
	googleCfg   *oauth2.Config
	githubCfg   *oauth2.Config
}

func NewOAuthHandler(
	authService *service.AuthService,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	errHandler *ErrorHandler,
) *OAuthHandler {
	h := &OAuthHandler{
		authService: authService,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		errHandler:  errHandler,
	}

	// Google OAuth2 config — set env vars GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET
	googleID := os.Getenv("GOOGLE_CLIENT_ID")
	googleSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleID != "" && googleSecret != "" {
		h.googleCfg = &oauth2.Config{
			ClientID:     googleID,
			ClientSecret: googleSecret,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
			RedirectURL: os.Getenv("OAUTH_REDIRECT_BASE") + "/auth/google/callback",
		}
	}

	// GitHub OAuth2 config — set env vars GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET
	githubID := os.Getenv("GITHUB_CLIENT_ID")
	githubSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if githubID != "" && githubSecret != "" {
		h.githubCfg = &oauth2.Config{
			ClientID:     githubID,
			ClientSecret: githubSecret,
			Scopes:       []string{"user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			RedirectURL: os.Getenv("OAUTH_REDIRECT_BASE") + "/auth/github/callback",
		}
	}

	return h
}

// Enabled returns true if at least one OAuth provider is configured
func (h *OAuthHandler) GoogleEnabled() bool { return h.googleCfg != nil }
func (h *OAuthHandler) GitHubEnabled() bool { return h.githubCfg != nil }

// generateState creates a random state token for CSRF protection
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ============================================
// GOOGLE
// ============================================

func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.googleCfg == nil {
		h.errHandler.Error(w, r, http.StatusServiceUnavailable, "Google OAuth non configuré")
		return
	}
	state := generateState()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, h.googleCfg.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if h.googleCfg == nil {
		h.errHandler.Error(w, r, http.StatusServiceUnavailable, "Google OAuth non configuré")
		return
	}

	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		h.errHandler.Error(w, r, http.StatusForbidden, "État OAuth invalide")
		return
	}
	// Clear state cookie
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Path: "/", MaxAge: -1})

	code := r.URL.Query().Get("code")
	if code == "" {
		h.errHandler.Error(w, r, http.StatusBadRequest, "Code d'autorisation manquant")
		return
	}

	token, err := h.googleCfg.Exchange(r.Context(), code)
	if err != nil {
		h.errHandler.Error(w, r, http.StatusInternalServerError, "Échec de l'échange de token")
		return
	}

	// Fetch user info from Google
	client := h.googleCfg.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var gUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &gUser); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	h.loginOrCreateOAuthUser(w, r, gUser.Email, gUser.Name, "google")
}

// ============================================
// GITHUB
// ============================================

func (h *OAuthHandler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	if h.githubCfg == nil {
		h.errHandler.Error(w, r, http.StatusServiceUnavailable, "GitHub OAuth non configuré")
		return
	}
	state := generateState()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, h.githubCfg.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	if h.githubCfg == nil {
		h.errHandler.Error(w, r, http.StatusServiceUnavailable, "GitHub OAuth non configuré")
		return
	}

	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		h.errHandler.Error(w, r, http.StatusForbidden, "État OAuth invalide")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Path: "/", MaxAge: -1})

	code := r.URL.Query().Get("code")
	if code == "" {
		h.errHandler.Error(w, r, http.StatusBadRequest, "Code d'autorisation manquant")
		return
	}

	token, err := h.githubCfg.Exchange(r.Context(), code)
	if err != nil {
		h.errHandler.Error(w, r, http.StatusInternalServerError, "Échec de l'échange de token")
		return
	}

	// Fetch user info from GitHub
	client := h.githubCfg.Client(r.Context(), token)

	// Get user profile
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var ghUser struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}
	json.Unmarshal(body, &ghUser)

	// GitHub can return empty email — fetch from emails endpoint
	if ghUser.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			emailBody, _ := io.ReadAll(emailResp.Body)
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			json.Unmarshal(emailBody, &emails)
			for _, e := range emails {
				if e.Primary && e.Verified {
					ghUser.Email = e.Email
					break
				}
			}
		}
	}

	if ghUser.Email == "" {
		h.errHandler.Error(w, r, http.StatusBadRequest, "Impossible de récupérer l'email GitHub")
		return
	}

	h.loginOrCreateOAuthUser(w, r, ghUser.Email, ghUser.Login, "github")
}

// ============================================
// SHARED: find or create user, create session
// ============================================

func (h *OAuthHandler) loginOrCreateOAuthUser(w http.ResponseWriter, r *http.Request, email, name, provider string) {
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)

	// Try to find existing user by email
	user, err := h.userRepo.FindByEmail(email)
	if err != nil {
		// User doesn't exist — create one
		// Generate a unique username if taken
		username := name
		if username == "" {
			username = strings.Split(email, "@")[0]
		}
		// Ensure unique username
		base := username
		for i := 1; h.userRepo.UsernameExists(username); i++ {
			username = fmt.Sprintf("%s_%d", base, i)
		}

		// OAuth users get a random impossible-to-guess password hash
		// They can never login with password — only via OAuth
		randomHash := "$2a$10$OAuthPlaceholder" + utils.GenerateUUID()

		user = &models.User{
			Email:        email,
			Username:     username,
			PasswordHash: randomHash,
			Role:         "user",
		}
		if createErr := h.userRepo.Create(user); createErr != nil {
			h.errHandler.InternalError(w, r, createErr)
			return
		}
	}

	// Invalidate old sessions and create a new one
	h.sessionRepo.DeleteByUserID(user.ID)
	session := &models.Session{
		ID:        utils.GenerateUUID(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := h.sessionRepo.Create(session); err != nil {
		h.errHandler.InternalError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Lax required for OAuth redirect
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
