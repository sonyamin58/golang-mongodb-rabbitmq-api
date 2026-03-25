package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	AccountNumber string        `gorm:"uniqueIndex;size:50;not null" json:"account_number"`
	AccountType  string         `gorm:"size:20;not null" json:"account_type"`
	Balance      float64        `gorm:"type:decimal(18,2);default:0" json:"balance"`
	Currency     string         `gorm:"size:3;default:USD" json:"currency"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type AccountRepository interface {
	Create(account *Account) error
	FindByID(id uint) (*Account, error)
	FindByAccountNumber(accountNumber string) (*Account, error)
	FindByUserID(userID uint) ([]Account, error)
	Update(account *Account) error
	Delete(id uint) error
	List(limit, offset int) ([]Account, error)
	Count() (int64, error)
}

func (Account) TableName() string {
	return "accounts"
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
