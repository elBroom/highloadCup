package schema

type RequestLocationVisits struct {
	FromDate *int32  `json:"fromDate"`
	ToDate   *int32  `json:"toDate"`
	FromAge  *uint8  `json:"fromAge"`
	ToAge    *uint8  `json:"toAge"`
	Gender   *string `json:"genser"`
}

type ResponceLocationVisits struct {
	Avg float32 `json:"avg"`
}
