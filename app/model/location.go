package model

type Location struct {
	ID       *uint32 `json:"id"`
	Place    *string `json:"place"`
	Country  *string `json:"country"`
	City     *string `json:"city"`
	Distance *uint32 `json:"distance"`
}
