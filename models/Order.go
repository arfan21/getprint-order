package models

import (
	validator "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Order struct {
	ID           uint          `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    null.Time     `gorm:"index" json:"deleted_at,omitempty"`
	PartnerID    uint          `gorm:"not null" json:"partner_id"`
	UserID       uint          `gorm:"not null" json:"user_id"`
	Status       string        `gorm:"default:'pending'" json:"status"`
	OrderDetails []OrderDetail `json:"order_details"`
}

type OrderDetail struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt null.Time `gorm:"index" json:"deleted_at,omitempty"`
	OrderID   uint      `gorm:"not null" json:"order_id"`
	File      string    `gorm:"not null" json:"file"`
	Path      string    `gorm:"not null" json:"path"`
}

func (o Order) Validate() error {
	return validator.Errors{
		"partner_id": validator.Validate(o.PartnerID, validator.Required),
		"user_id":    validator.Validate(o.UserID, validator.Required),
	}.Filter()
}

type OrderRepository interface {
	Create(data *Order) error
	GetByID(id uint) (*Order, error)
	GetByUserID(userID uint) (*[]Order, error)
	Update(data *Order) error
}

type OrderService interface {
	Create(data *Order) error
	GetByID(id uint) (*Order, error)
	GetByUserID(userID uint) (*[]Order, error)
	Update(id uint, data *Order) error
}
