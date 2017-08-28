package model

type Visit struct {
	ID         *uint32 `json:"id,omitempty"`
	LocationID *uint32 `json:"location"`
	UserID     *uint32 `json:"user"`
	VisitedAt  *int32  `json:"visited_at"`
	Mark       *uint8  `json:"mark"`

	User     *User
	Location *Location
}

type Visits struct {
	Visits []Visit `json:"visits"`
}
