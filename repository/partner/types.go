package partner

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4/zero"
)

type PartnerResponse struct {
	Status  string          `json:"status"`
	Message interface{}     `json:"message,omitempty"`
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
	Print         zero.Int  `json:"print"`
	Scan          zero.Int  `json:"scan"`
	Fotocopy      zero.Int  `json:"fotocopy"`
	TotalFollower int64     `json:"total_follower"`
}
