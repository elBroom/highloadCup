package handler

import (
	"net/http"
	"strconv"

	"strings"

	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func writeObj(ctx *fasthttp.RequestCtx, obj easyjson.Marshaler) {
	b, err := easyjson.Marshal(obj)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetBody(b)
}

func writeStr(ctx *fasthttp.RequestCtx, s string) {
	b := []byte(s)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.SetBody(b)
}

func CheckNull(b []byte) bool {
	return string(b) == "" || strings.Contains(string(b), ": null")
}
