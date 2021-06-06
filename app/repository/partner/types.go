package partner

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type PartnerResponse struct {
	Status  string        `json:"status"`
	Message interface{}   `json:"message,omitempty"`
	Data    partnerStruct `json:"data,omitempty"`
}

type partnerStruct struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        uuid.UUID `json:"user_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Picture       string    `json:"picture"`
	Address       string    `json:"address"`
	Lat           string    `json:"lat"`
	Lng           string    `json:"lng"`
	Print         int64     `json:"print"`
	Scan          int64     `json:"scan"`
	Fotocopy      int64     `json:"fotocopy"`
	TotalFollower int64     `json:"total_follower"`
}
