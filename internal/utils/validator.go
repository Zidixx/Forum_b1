package utils

import (
	"net/mail"
	"strings"
)

type ValidationErrors map[string]string

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

func ValidateRegister(email, username, password, confirmPassword string) ValidationErrors {
	errs := make(ValidationErrors)

	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)

	if email == "" {
		errs["email"] = "L'email est requis"
	} else if _, err := mail.ParseAddress(email); err != nil {
		errs["email"] = "Email invalide"
	}

	if username == "" {
		errs["username"] = "Le nom d'utilisateur est requis"
	} else if len(username) < 3 {
		errs["username"] = "Le nom d'utilisateur doit contenir au moins 3 caractères"
	} else if len(username) > 30 {
		errs["username"] = "Le nom d'utilisateur ne doit pas dépasser 30 caractères"
	}

	if password == "" {
		errs["password"] = "Le mot de passe est requis"
	} else if len(password) < 6 {
		errs["password"] = "Le mot de passe doit contenir au moins 6 caractères"
	}

	if confirmPassword != password {
		errs["confirm_password"] = "Les mots de passe ne correspondent pas"
	}

	return errs
}

func ValidateLogin(identifier, password string) ValidationErrors {
	errs := make(ValidationErrors)

	if strings.TrimSpace(identifier) == "" {
		errs["identifier"] = "L'email ou le nom d'utilisateur est requis"
	}
	if password == "" {
		errs["password"] = "Le mot de passe est requis"
	}

	return errs
}

func ValidatePost(title, content string, categoryIDs []int) ValidationErrors {
	errs := make(ValidationErrors)

	if strings.TrimSpace(title) == "" {
		errs["title"] = "Le titre est requis"
	} else if len(title) > 200 {
		errs["title"] = "Le titre ne doit pas dépasser 200 caractères"
	}

	if strings.TrimSpace(content) == "" {
		errs["content"] = "Le contenu est requis"
	}

	if len(categoryIDs) == 0 {
		errs["categories"] = "Au moins une catégorie est requise"
	}

	return errs
}

func ValidateComment(content string) ValidationErrors {
	errs := make(ValidationErrors)

	if strings.TrimSpace(content) == "" {
		errs["content"] = "Le commentaire ne peut pas être vide"
	}

	return errs
}
