package rassilki

import (
	"errors"
	"github.com/lib/pq"
	"time"
)

type Rassilka struct {
	Id           int            `json:"-"              db:"id"`
	StartTime    time.Time      `json:"start-time"     db:"start_time"   binding:"required"`
	Message      string         `json:"message"        db:"message"      binding:"required"`
	Filter       pq.StringArray `json:"filter"         db:"filter"`
	EndTime      time.Time      `json:"end-time"       db:"end_time"     binding:"required"`
	Supplemented bool           `json:"-"              db:"supplemented"`
	Recreated    bool           `json:"-"              db:"recreated"`
}

type UpdateRassilka struct {
	StartTime    *time.Time `json:"start-time" `
	Message      *string    `json:"message" `
	Filter       *[]string  `json:"filter"`
	EndTime      *time.Time `json:"end-time" `
	Supplemented *bool      `json:"supplemented"`
	Recreated    *bool      `json:"recreated"`
}

func (i UpdateRassilka) Validate() error {
	if i.StartTime == nil && i.Message == nil && i.EndTime == nil {
		return errors.New("updating structure has no values")
	}
	if i.StartTime != nil && i.EndTime != nil {
		if (*i.StartTime).After(*i.EndTime) {
			return errors.New("StartTime is after EndTime")
		}
	}
	if i.EndTime != nil && (*i.EndTime).Before(time.Now()) {
		return errors.New("EndTime is in the past")
	}
	return nil
}

type RassilkaReview struct {
	Id      int `json:"ras_id"         db:"ras_id"`
	Total   int `json:"total"          db:"-"`
	Sent    int `json:"sent"           db:"-"`
	NotSent int `json:"not_sent"       db:"-"`
}
