package model

import "strings"

type Visit struct {
	ID         *uint32 `json:"id,omitempty"`
	LocationID *uint32 `json:"location"`
	UserID     *uint32 `json:"user"`
	VisitedAt  *int64  `json:"visited_at"`
	Mark       *uint8  `json:"mark"`
}

type Visits struct {
	Visits []Visit `json:"visits"`
}

func (v *Visits) CheckNull(b []byte) bool {
	return strings.Contains(string(b), ": null")
}
