package service

import (
	"errors"
	"time"

	"github.com/ibas/golib-api/internal/config"
	"github.com/ibas/golib-api/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo model.UserRepository
	cfg      *config.Config
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func NewAuthService(userRepo model.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, err = s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     "user",
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*TokenPair, *model.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, nil, errors.New("account is inactive")
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	user.Password = ""
	return tokens, user, nil
}

func (s *AuthService) RefreshToken(req *RefreshRequest) (*TokenPair, error) {
	claims, err := s.validateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user, err := s.userRepo.FindByID(uint(userID))
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsActive {
		return nil, errors.New("user not found or inactive")
	}

	return s.generateTokens(user)
}

func (s *AuthService) generateTokens(user *model.User) (*TokenPair, error) {
	now := time.Now()

	// Access token
	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"iat":      now.Unix(),
		"exp":      now.Add(s.cfg.JWT.GetExpiryDuration()).Unix(),
		"type":     "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"iat":      now.Unix(),
		"exp":      now.Add(s.cfg.JWT.GetRefreshExpiryDuration()).Unix(),
		"type":     "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:   int64(s.cfg.JWT.GetExpiryDuration().Seconds()),
	}, nil
}

func (s *AuthService) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	return s.validateToken(tokenString)
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
