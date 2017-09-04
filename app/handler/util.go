package handler

import (
	"net/http"
	"strconv"

	"math"

	"strings"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
)

func writeObj(ctx *fasthttp.RequestCtx, obj interface{}) {
	b, err := ffjson.Marshal(obj)
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

func parseID(ctx *fasthttp.RequestCtx) (uint32, error) {
	id, err := strconv.ParseUint(string(ctx.UserValue("id").(string)), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func CheckNull(b []byte) bool {
	return string(b) == "" || strings.Contains(string(b), ": null")
}
