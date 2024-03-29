package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type Cart struct {
	ID        uint        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt null.Time   `gorm:"index" json:"deleted_at,omitempty"`
	UserID    uuid.UUID   `gorm:"type:char(36);not null" json:"user_id"`
	PartnerID uint        `gorm:"not null" json:"partner_id"`
	OrderType string      `gorm:"type:enum('print','fotocopy','scan');not null" json:"order_type"`
	Qty       uint        `gorm:"default:1" json:"qty"`
	LinkFile  null.String `json:"link_file"`
}
