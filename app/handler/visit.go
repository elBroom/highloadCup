package handler

import (
	"net/http"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetVisitEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	visit, ok := storage.DataStorage.Visit.Get(id)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	writeObj(ctx, visit)
}

func UpdateVisitEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	bytes := ctx.PostBody()
	ok := CheckNull(bytes)
	if ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	var visit model.Visit
	_ = easyjson.Unmarshal(bytes, &visit)

	err := storage.DataStorage.Visit.Update(id, &visit)
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
	_ = easyjson.Unmarshal(ctx.PostBody(), &visit)

	err := storage.DataStorage.Visit.Add(&visit)
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
