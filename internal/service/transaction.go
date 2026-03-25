package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ibas/golib-api/internal/machinery"
	"github.com/ibas/golib-api/internal/model"
)

type TransactionService struct {
	transactionRepo model.TransactionRepository
	accountRepo     model.AccountRepository
	machineryClient *machinery.AsyncClient
}

func NewTransactionService(transactionRepo model.TransactionRepository, accountRepo model.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func NewTransactionServiceWithMachinery(transactionRepo model.TransactionRepository, accountRepo model.AccountRepository, client *machinery.AsyncClient) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		machineryClient: client,
	}
}

type CreateTransactionRequest struct {
	AccountID       uint   `json:"account_id"`
	TransactionType string `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	Description     string `json:"description"`
}

type TransferRequest struct {
	FromAccountID uint    `json:"from_account_id"`
	ToAccountID   uint    `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
}

func (s *TransactionService) Create(req *CreateTransactionRequest) (*model.Transaction, error) {
	// Verify account exists
	account, err := s.accountRepo.FindByID(req.AccountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	if !account.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Validate transaction type
	txType := model.TransactionType(req.TransactionType)
	if txType != model.TransactionTypeDeposit && txType != model.TransactionTypeWithdrawal {
		return nil, errors.New("invalid transaction type")
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Check sufficient balance for withdrawals
	if txType == model.TransactionTypeWithdrawal && account.Balance < req.Amount {
		return nil, errors.New("insufficient balance")
	}

	// Generate reference number
	refNumber, err := s.generateReferenceNumber()
	if err != nil {
		return nil, err
	}

	// Create transaction
	balanceBefore := account.Balance
	var balanceAfter float64

	if txType == model.TransactionTypeDeposit {
		balanceAfter = account.Balance + req.Amount
	} else {
		balanceAfter = account.Balance - req.Amount
	}

	transaction := &model.Transaction{
		AccountID:       req.AccountID,
		TransactionType: txType,
		Amount:          req.Amount,
		Currency:        account.Currency,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Description:     req.Description,
		ReferenceNumber: refNumber,
		Status:          model.TransactionStatusCompleted,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	// Update account balance
	account.Balance = balanceAfter
	if err := s.accountRepo.Update(account); err != nil {
		return nil, err
	}

	// Send async task if Machinery client is configured
	if s.machineryClient != nil {
		go func() {
			var taskErr error
			if txType == model.TransactionTypeDeposit {
				_, taskErr = s.machineryClient.SendTopupTask(transaction.ID)
			} else {
				_, taskErr = s.machineryClient.SendWithdrawTask(transaction.ID)
			}
			if taskErr != nil {
				fmt.Printf("Failed to send async task: %v\n", taskErr)
			}
		}()
	}

	return transaction, nil
}

func (s *TransactionService) Transfer(req *TransferRequest) (*model.Transaction, *model.Transaction, error) {
	// Verify source account
	fromAccount, err := s.accountRepo.FindByID(req.FromAccountID)
	if err != nil {
		return nil, nil, err
	}
	if fromAccount == nil {
		return nil, nil, errors.New("source account not found")
	}
	if !fromAccount.IsActive {
		return nil, nil, errors.New("source account is inactive")
	}

	// Verify destination account
	toAccount, err := s.accountRepo.FindByID(req.ToAccountID)
	if err != nil {
		return nil, nil, err
	}
	if toAccount == nil {
		return nil, nil, errors.New("destination account not found")
	}
	if !toAccount.IsActive {
		return nil, nil, errors.New("destination account is inactive")
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, nil, errors.New("amount must be greater than zero")
	}

	// Check sufficient balance
	if fromAccount.Balance < req.Amount {
		return nil, nil, errors.New("insufficient balance")
	}

	// Generate reference numbers
	refNumber, err := s.generateReferenceNumber()
	if err != nil {
		return nil, nil, err
	}

	// Create withdrawal transaction
	fromBalanceBefore := fromAccount.Balance
	fromBalanceAfter := fromAccount.Balance - req.Amount

	fromTransaction := &model.Transaction{
		AccountID:        req.FromAccountID,
		RelatedAccountID: &req.ToAccountID,
		TransactionType:  model.TransactionTypeTransfer,
		Amount:           req.Amount,
		Currency:         fromAccount.Currency,
		BalanceBefore:    fromBalanceBefore,
		BalanceAfter:     fromBalanceAfter,
		Description:      req.Description,
		ReferenceNumber:  refNumber + "-OUT",
		Status:           model.TransactionStatusCompleted,
	}

	if err := s.transactionRepo.Create(fromTransaction); err != nil {
		return nil, nil, err
	}

	// Create deposit transaction
	toBalanceBefore := toAccount.Balance
	toBalanceAfter := toAccount.Balance + req.Amount

	toTransaction := &model.Transaction{
		AccountID:        req.ToAccountID,
		RelatedAccountID: &req.FromAccountID,
		TransactionType:  model.TransactionTypeTransfer,
		Amount:           req.Amount,
		Currency:         toAccount.Currency,
		BalanceBefore:    toBalanceBefore,
		BalanceAfter:     toBalanceAfter,
		Description:      req.Description,
		ReferenceNumber:  refNumber + "-IN",
		Status:           model.TransactionStatusCompleted,
	}

	if err := s.transactionRepo.Create(toTransaction); err != nil {
		return nil, nil, err
	}

	// Update account balances
	fromAccount.Balance = fromBalanceAfter
	toAccount.Balance = toBalanceAfter

	if err := s.accountRepo.Update(fromAccount); err != nil {
		return nil, nil, err
	}
	if err := s.accountRepo.Update(toAccount); err != nil {
		return nil, nil, err
	}

	// Send async task if Machinery client is configured
	if s.machineryClient != nil {
		go func() {
			_, err := s.machineryClient.SendTransferTask(fromTransaction.ID)
			if err != nil {
				fmt.Printf("Failed to send transfer task: %v\n", err)
			}
		}()
	}

	return fromTransaction, toTransaction, nil
}

func (s *TransactionService) Get(id uint) (*model.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errors.New("transaction not found")
	}
	return transaction, nil
}

func (s *TransactionService) GetByReferenceNumber(refNumber string) (*model.Transaction, error) {
	return s.transactionRepo.FindByReferenceNumber(refNumber)
}

func (s *TransactionService) GetByAccountID(accountID uint, limit, offset int) ([]model.Transaction, error) {
	return s.transactionRepo.FindByAccountID(accountID, limit, offset)
}

func (s *TransactionService) List(limit, offset int) ([]model.Transaction, int64, error) {
	transactions, err := s.transactionRepo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.transactionRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

func (s *TransactionService) generateReferenceNumber() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("TXN%s%d", hex.EncodeToString(bytes)[:8], time.Now().Unix()), nil
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
