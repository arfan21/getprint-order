package media

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type MediaResponse struct{
	Status  string        `json:"status"`
	Message interface{}   `json:"message,omitempty"`
	Data    mediaResponse `json:"data,omitempty"`
}

type mediaResponse struct{

	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt null.Time `gorm:"index" json:"deleted_at,omitempty"`
	Url       string    `json:"url"`
	Path      string    `json:"path"`
}