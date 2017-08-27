package schema

import "time"

type RequestLocationVisits struct {
	FromDate *time.Time `json:"fromDate"`
	ToDate   *time.Time `json:"toDate"`
	FromAge  *uint8     `json:"fromAge"`
	ToAge    *uint8     `json:"toAge"`
	Gender   *string    `json:"genser"`
}

type ResponceLocationVisits struct {
	Age *float32 `json:"age"`
}
