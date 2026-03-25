package handler

import (
	"net/http"
	"strconv"

	"github.com/ibas/golib-api/internal/response"
	"github.com/ibas/golib-api/internal/service"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) List(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	transactions, total, err := h.transactionService.List(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Pagination(transactions, total, limit, offset))
}

func (h *TransactionHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid transaction ID"))
	}

	transaction, err := h.transactionService.Get(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(transaction))
}

func (h *TransactionHandler) Create(c echo.Context) error {
	var req service.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
	}

	transaction, err := h.transactionService.Create(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(transaction))
}

func (h *TransactionHandler) Transfer(c echo.Context) error {
	var req service.TransferRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
	}

	fromTx, toTx, err := h.transactionService.Transfer(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(map[string]interface{}{
		"from_transaction": fromTx,
		"to_transaction":   toTx,
	}))
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
