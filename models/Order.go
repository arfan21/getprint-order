package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type Order struct {
	ID           uint          `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    null.Time     `gorm:"index" json:"deleted_at,omitempty"`
	PartnerID    uint          `gorm:"not null" json:"partner_id"`
	UserID       uuid.UUID     `gorm:"not null" json:"user_id"`
	TotalPrice   uint          `gorm:"not null" json:"total_price"`
	Status       string        `gorm:"default:'pending'" json:"status"`
	OrderDetails []OrderDetail `gorm:"foreignKey:order_id;constraint:OnDelete:CASCADE;" json:"order_details"`
}

type OrderDetail struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    null.Time `gorm:"index" json:"deleted_at,omitempty"`
	OrderID      uint      `gorm:"not null" json:"order_id"`
	File         string    `gorm:"not null" json:"file"`
	Path         string    `gorm:"not null" json:"path"`
	PrintQty     uint      `gorm:"default:0" json:"print_qty"`
	ScanQty      uint      `gorm:"default:0" json:"scan_qty"`
	PhotocopyQty uint      `gorm:"default:0" json:"photocopy_qty"`
}
