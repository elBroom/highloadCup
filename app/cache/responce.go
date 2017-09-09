package cache

import (
	"time"
)

type Responce struct {
	life          time.Time
	statusCode    int
	body          []byte
	contentLength int
}

func (r *Responce) StatusCode() *int {
	return &r.statusCode
}

func (r *Responce) Body() *[]byte {
	return &r.body
}
