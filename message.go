package rassilki

import (
	"time"
)

type Message struct {
	Id       int        `json:"id"        db:"id"        binding:"required"`
	SentAt   *time.Time `json:"sent-at"   db:"sent_at"   binding:"required"`
	Status   bool       `json:"status"    db:"status"    binding:"required"`
	RasId    int        `json:"ras-id"    db:"ras_id"    binding:"required"`
	ClientId int        `json:"client-id" db:"client_id" binding:"required""`
}
