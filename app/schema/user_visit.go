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

func (r *ResponceUserVisits) Len() int {
	return len(r.Visits)
}

func (r *ResponceUserVisits) Less(i, j int) bool {
	return (*(*r.Visits[i]).Visited_at) < (*(*r.Visits[j]).Visited_at)
}

func (r *ResponceUserVisits) Swap(i, j int) {
	r.Visits[i], r.Visits[j] = r.Visits[j], r.Visits[i]
}
