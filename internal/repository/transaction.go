package repository

import (
	"errors"

	"github.com/ibas/golib-api/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) FindByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Preload("Account").First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) FindByReferenceNumber(refNumber string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Where("reference_number = ?", refNumber).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) FindByAccountID(accountID uint, limit, offset int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("account_id = ?", accountID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Update(transaction *model.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *TransactionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Transaction{}, id).Error
}

func (r *TransactionRepository) List(limit, offset int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("Account").Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Transaction{}).Count(&count).Error
	return count, err
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
