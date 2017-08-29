package model

type Visit struct {
	ID         *uint32 `json:"id,omitempty"`
	LocationID *uint32 `json:"location"`
	UserID     *uint32 `json:"user"`
	VisitedAt  *int64  `json:"visited_at"`
	Mark       *uint8  `json:"mark"`

	User     *User     `json:"-"`
	Location *Location `json:"-"`
}

type Visits struct {
	Visits []Visit `json:"visits"`
}
