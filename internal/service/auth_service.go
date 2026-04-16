package service

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/utils"
	"strings"
	"time"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *AuthService) Register(email, username, password string) (*models.User, utils.ValidationErrors) {
	errs := make(utils.ValidationErrors)

	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)

	if s.userRepo.EmailExists(email) {
		errs["email"] = "Cet email est déjà utilisé"
	}
	if s.userRepo.UsernameExists(username) {
		errs["username"] = "Ce nom d'utilisateur est déjà pris"
	}
	if errs.HasErrors() {
		return nil, errs
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		errs["general"] = "Erreur interne"
		return nil, errs
	}

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		errs["general"] = "Erreur lors de la création du compte"
		return nil, errs
	}

	return user, nil
}

func (s *AuthService) Login(identifier, password string) (*models.User, *models.Session, error) {
	identifier = strings.TrimSpace(identifier)

	user, err := s.userRepo.FindByEmail(identifier)
	if err != nil {
		user, err = s.userRepo.FindByUsername(identifier)
		if err != nil {
			return nil, nil, fmt.Errorf("identifiants incorrects")
		}
	}

	if !utils.CheckPassword(user.PasswordHash, password) {
		return nil, nil, fmt.Errorf("identifiants incorrects")
	}

	// Invalidate existing sessions
	s.sessionRepo.DeleteByUserID(user.ID)

	session := &models.Session{
		ID:        utils.GenerateUUID(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, nil, fmt.Errorf("erreur création de session")
	}

	return user, session, nil
}

func (s *AuthService) Logout(sessionID string) {
	s.sessionRepo.DeleteByID(sessionID)
}

func (s *AuthService) GetUserFromSession(sessionID string) (*models.User, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.sessionRepo.DeleteByID(sessionID)
		return nil, fmt.Errorf("session expirée")
	}

	return s.userRepo.FindByID(session.UserID)
}
