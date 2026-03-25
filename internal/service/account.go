package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ibas/golib-api/internal/model"
)

type AccountService struct {
	accountRepo model.AccountRepository
	userRepo    model.UserRepository
}

func NewAccountService(accountRepo model.AccountRepository, userRepo model.UserRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
	}
}

type CreateAccountRequest struct {
	UserID      uint    `json:"user_id"`
	AccountType string  `json:"account_type"`
	InitialBalance float64 `json:"initial_balance"`
	Currency    string  `json:"currency"`
}

type UpdateAccountRequest struct {
	IsActive    *bool   `json:"is_active"`
	Currency    string  `json:"currency"`
}

func (s *AccountService) Create(req *CreateAccountRequest) (*model.Account, error) {
	// Verify user exists
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Generate unique account number
	accountNumber, err := s.generateAccountNumber()
	if err != nil {
		return nil, err
	}

	account := &model.Account{
		UserID:        req.UserID,
		AccountNumber: accountNumber,
		AccountType:   req.AccountType,
		Balance:       req.InitialBalance,
		Currency:      req.Currency,
		IsActive:      true,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) Get(id uint) (*model.Account, error) {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (s *AccountService) GetByAccountNumber(accountNumber string) (*model.Account, error) {
	account, err := s.accountRepo.FindByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (s *AccountService) GetByUserID(userID uint) ([]model.Account, error) {
	return s.accountRepo.FindByUserID(userID)
}

func (s *AccountService) Update(id uint, req *UpdateAccountRequest) (*model.Account, error) {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}

	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}
	if req.Currency != "" {
		account.Currency = req.Currency
	}

	if err := s.accountRepo.Update(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) Delete(id uint) error {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return err
	}
	if account == nil {
		return errors.New("account not found")
	}

	return s.accountRepo.Delete(id)
}

func (s *AccountService) List(limit, offset int) ([]model.Account, int64, error) {
	accounts, err := s.accountRepo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.accountRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return accounts, count, nil
}

func (s *AccountService) generateAccountNumber() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("ACC%s", hex.EncodeToString(bytes)[:12]), nil
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
