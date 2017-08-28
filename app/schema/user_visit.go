package schema

type RequestUserVisits struct {
	FromDate   *int32  `json:"fromDate"`
	ToDate     *int32  `json:"toDate"`
	Country    *string `json:"country"`
	ToDistance *uint32 `json:"toDistance"`
}

type ResponceUserVisit struct {
	Mark       *uint8  `json:"mark"`
	Visited_at *int32  `json:"visited_at"`
	Place      *string `json:"place"`
}

type ResponceUserVisits struct {
	Visits []*ResponceUserVisit `json:"visits"`
}
