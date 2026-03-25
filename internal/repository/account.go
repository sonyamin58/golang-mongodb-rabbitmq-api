package repository

import (
	"errors"

	"github.com/ibas/golib-api/internal/model"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *AccountRepository) FindByID(id uint) (*model.Account, error) {
	var account model.Account
	err := r.db.Preload("User").First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindByAccountNumber(accountNumber string) (*model.Account, error) {
	var account model.Account
	err := r.db.Preload("User").Where("account_number = ?", accountNumber).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindByUserID(userID uint) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) Update(account *model.Account) error {
	return r.db.Save(account).Error
}

func (r *AccountRepository) Delete(id uint) error {
	return r.db.Delete(&model.Account{}, id).Error
}

func (r *AccountRepository) List(limit, offset int) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Preload("User").Limit(limit).Offset(offset).Order("created_at DESC").Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Account{}).Count(&count).Error
	return count, err
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
