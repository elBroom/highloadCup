package cache

import (
	"time"

	"sync"

	"github.com/valyala/fasthttp"
)

type Memory struct {
	mx   sync.RWMutex
	data map[string]*Responce
}

func key(uri *fasthttp.URI) string {
	args := uri.QueryArgs()
	args.Del("query_id")

	return string(append(uri.Path(), args.QueryString()...))
}

func (m *Memory) Get(uri *fasthttp.URI) (*Responce, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	resp, ok := m.data[key(uri)]
	if ok {
		if !time.Now().After(resp.life) {
			return resp, ok
		}
	}
	return nil, false
}

func (m *Memory) Set(uri *fasthttp.URI, responce *fasthttp.Response, lifeTime time.Duration) {
	m.mx.Lock()
	defer m.mx.Unlock()

	body := string(responce.Body())
	m.data[key(uri)] = &Responce{
		life:       time.Now().Add(lifeTime),
		statusCode: responce.StatusCode(),
		body:       []byte(body),
	}
}

func InitCache() *Memory {
	return &Memory{data: make(map[string]*Responce)}
}
