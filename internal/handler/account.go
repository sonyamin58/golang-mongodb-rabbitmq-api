package handler

import (
	"net/http"
	"strconv"

	"github.com/ibas/golib-api/internal/response"
	"github.com/ibas/golib-api/internal/service"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) List(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	accounts, total, err := h.accountService.List(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Pagination(accounts, total, limit, offset))
}

func (h *AccountHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid account ID"))
	}

	account, err := h.accountService.Get(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(account))
}

func (h *AccountHandler) Create(c echo.Context) error {
	var req service.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
	}

	if req.AccountType == "" {
		req.AccountType = "CHECKING"
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}
	if req.InitialBalance < 0 {
		return c.JSON(http.StatusBadRequest, response.Error("Initial balance cannot be negative"))
	}

	account, err := h.accountService.Create(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(account))
}

func (h *AccountHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid account ID"))
	}

	var req service.UpdateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
	}

	account, err := h.accountService.Update(uint(id), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(account))
}

func (h *AccountHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid account ID"))
	}

	if err := h.accountService.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(map[string]string{
		"message": "Account deleted successfully",
	}))
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
