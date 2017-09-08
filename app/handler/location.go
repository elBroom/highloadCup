package handler

import (
	"net/http"

	"fmt"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetLocationEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	location, ok := storage.DataStorage.Location.Get(id)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	writeObj(ctx, location)
}

func GetLocatioAvgEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	params := ctx.QueryArgs()

	//  Parse fromDate parameter
	var fromDate *int64
	if params.Has("fromDate") {
		if _, err := params.GetUint("fromDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("fromDate")
		tmp2 := int64(tmp)
		fromDate = &tmp2
	}

	//  Parse toDate parameter
	var toDate *int64
	if params.Has("toDate") {
		if _, err := params.GetUint("toDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toDate")
		tmp2 := int64(tmp)
		toDate = &tmp2
	}

	//  Parse fromAge parameter
	var fromAge *int64
	if params.Has("fromAge") {
		if _, err := params.GetUint("fromAge"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("fromAge")
		tmp2 := int64(tmp)
		fromAge = &tmp2
	}

	//  Parse toAge parameter
	var toAge *int64
	if params.Has("toAge") {
		if _, err := params.GetUint("toAge"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toAge")
		tmp2 := int64(tmp)
		toAge = &tmp2
	}

	//  Parse gender parameter
	var gender *string
	if params.Has("gender") {
		tmp := string(params.Peek("gender"))
		if tmp != "m" && tmp != "f" {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		gender = &tmp
	}

	avg := storage.DataStorage.VisitList.GetByLocation(id, fromDate, toDate, fromAge, toAge, gender)
	if avg == nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}
	answ := fmt.Sprintf(`{"avg": %.5f}`, *avg)
	writeStr(ctx, answ)
}

func UpdateLocationEndpoint(ctx *fasthttp.RequestCtx, id uint32) {
	bytes := ctx.PostBody()
	ok := CheckNull(bytes)
	if ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	var location model.Location
	_ = easyjson.Unmarshal(bytes, &location)

	err := storage.DataStorage.Location.Update(id, &location)
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

func CreateLocationEndpoint(ctx *fasthttp.RequestCtx) {
	var location model.Location
	_ = easyjson.Unmarshal(ctx.PostBody(), &location)

	err := storage.DataStorage.Location.Add(&location)
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
