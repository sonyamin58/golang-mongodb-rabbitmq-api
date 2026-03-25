package validator

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailInvalid     = errors.New("invalid email format")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordWeak     = errors.New("password must be at least 8 characters")
	ErrUsernameRequired = errors.New("username is required")
	ErrUsernameInvalid  = errors.New("username must be alphanumeric")
)

type Validator struct {
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
}

func New() *Validator {
	return &Validator{
		emailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
	}
}

func (v *Validator) ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrEmailRequired
	}
	if !v.emailRegex.MatchString(email) {
		return ErrEmailInvalid
	}
	return nil
}

func (v *Validator) ValidatePassword(password string) error {
	password = strings.TrimSpace(password)
	if password == "" {
		return ErrPasswordRequired
	}
	if len(password) < 8 {
		return ErrPasswordWeak
	}
	return nil
}

func (v *Validator) ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return ErrUsernameRequired
	}
	if len(username) < 3 || len(username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !v.usernameRegex.MatchString(username) {
		return ErrUsernameInvalid
	}
	return nil
}

func (v *Validator) ValidatePasswordStrength(password string) (bool, string) {
	if len(password) < 8 {
		return false, "password must be at least 8 characters"
	}
	
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}
	
	score := 0
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}
	
	if score < 3 {
		return false, "password must contain at least 3 of: uppercase, lowercase, digits, special characters"
	}
	
	return true, ""
}

func ValidateEmail(email string) error {
	return New().ValidateEmail(email)
}

func ValidatePassword(password string) error {
	return New().ValidatePassword(password)
}

func ValidateUsername(username string) error {
	return New().ValidateUsername(username)
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
