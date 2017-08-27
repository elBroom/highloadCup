package schema

import "time"

type RequestUserVisits struct {
	FromDate   *time.Time `json:"fromDate"`
	ToDate     *time.Time `json:"toDate"`
	Country    *string    `json:"country"`
	ToDistance *uint32    `json:"toDistance"`
}

type ResponceUserVisits struct {
	Mark       *uint8     `json:"mark"`
	Visited_at *time.Time `json:"visited_at"`
	Place      *string    `json:"place"`
}
