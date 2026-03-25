package handler

import (
	"net/http"

	"github.com/ibas/golib-api/internal/response"
	"github.com/ibas/golib-api/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req service.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
	}

	// Validate required fields
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, response.Error("Email, username, and password are required"))
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(user))
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req service.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, response.Error("Email and password are required"))
	}

	tokens, user, err := h.authService.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(map[string]interface{}{
		"tokens": tokens,
		"user":  user,
	}))
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req service.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, response.Error("Refresh token is required"))
	}

	tokens, err := h.authService.RefreshToken(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(tokens))
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
