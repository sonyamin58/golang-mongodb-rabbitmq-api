package model

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal  TransactionType = "WITHDRAWAL"
	TransactionTypeTransfer    TransactionType = "TRANSFER"
	TransactionTypePayment     TransactionType = "PAYMENT"
	
	TransactionStatusPending  TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed   TransactionStatus = "FAILED"
	TransactionStatusReversed TransactionStatus = "REVERSED"
)

type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	AccountID       uint              `gorm:"index;not null" json:"account_id"`
	RelatedAccountID *uint            `gorm:"index" json:"related_account_id,omitempty"`
	TransactionType TransactionType   `gorm:"size:20;not null" json:"transaction_type"`
	Amount          float64           `gorm:"type:decimal(18,2);not null" json:"amount"`
	Currency        string            `gorm:"size:3;default:USD" json:"currency"`
	BalanceBefore   float64           `gorm:"type:decimal(18,2)" json:"balance_before"`
	BalanceAfter    float64           `gorm:"type:decimal(18,2)" json:"balance_after"`
	Description     string            `gorm:"type:text" json:"description"`
	ReferenceNumber string            `gorm:"size:100;uniqueIndex" json:"reference_number"`
	Status          TransactionStatus `gorm:"size:20;default:PENDING" json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
	
	Account         Account           `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

type TransactionRepository interface {
	Create(transaction *Transaction) error
	FindByID(id uint) (*Transaction, error)
	FindByReferenceNumber(refNumber string) (*Transaction, error)
	FindByAccountID(accountID uint, limit, offset int) ([]Transaction, error)
	Update(transaction *Transaction) error
	Delete(id uint) error
	List(limit, offset int) ([]Transaction, error)
	Count() (int64, error)
}

func (Transaction) TableName() string {
	return "transactions"
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
