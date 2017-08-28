package schema

type RequestLocationVisits struct {
	FromDate *int64  `json:"fromDate"`
	ToDate   *int64  `json:"toDate"`
	FromAge  *int    `json:"fromAge"`
	ToAge    *int    `json:"toAge"`
	Gender   *string `json:"genser"`
}

type ResponceLocationVisits struct {
	Avg float64 `json:"avg"`
}
