package schema

type RequestUserVisits struct {
	FromDate   *int64  `json:"fromDate"`
	ToDate     *int64  `json:"toDate"`
	Country    *string `json:"country"`
	ToDistance *uint32 `json:"toDistance"`
}

type ResponceUserVisit struct {
	Mark       *uint8  `json:"mark"`
	Visited_at *int64  `json:"visited_at"`
	Place      *string `json:"place"`
}

type ResponceUserVisits struct {
	Visits []*ResponceUserVisit `json:"visits"`
}
