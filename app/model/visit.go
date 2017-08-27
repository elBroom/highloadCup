package model

import (
	"time"
)

type Visit struct {
	ID         *uint32    `json:"id,omitempty"`
	LocationID *uint32    `json:"location"`
	UserID     *uint32    `json:"user"`
	VisitedAt  *time.Time `json:"visited_at"`
	Mark       *uint8     `json:"mark"`

	User     *User
	Location *Location
}
