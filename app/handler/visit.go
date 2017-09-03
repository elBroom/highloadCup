package handler

import (
	"net/http"

	"encoding/json"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/valyala/fasthttp"
)

func GetVisitEndpoint(ctx *fasthttp.RequestCtx) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	visit, ok := storage.DataStorage.Visit.Get(id)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	writeObj(ctx, visit)
}

func UpdateVisitEndpoint(ctx *fasthttp.RequestCtx) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	bytes := ctx.PostBody()
	ok := CheckNull(bytes)
	if ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	var visit model.Visit
	_ = json.Unmarshal(bytes, &visit)

	err = storage.DataStorage.Visit.Update(id, &visit, storage.DataStorage)
	if err != nil {

		if err == storage.ErrDoesNotExist {
			ctx.SetStatusCode(http.StatusNotFound)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		return
	}

	writeStr(ctx, "{}")
}

func CreateVisitEndpoint(ctx *fasthttp.RequestCtx) {
	var visit model.Visit
	_ = json.Unmarshal(ctx.PostBody(), &visit)

	err := storage.DataStorage.Visit.Add(&visit, storage.DataStorage)
	if err != nil {

		if err == storage.ErrDoesNotExist {
			ctx.SetStatusCode(http.StatusNotFound)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		return
	}

	writeStr(ctx, "{}")
}
