package handler

import (
	"net/http"

	"fmt"

	"time"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
)

func GetLocationEndpoint(ctx *fasthttp.RequestCtx) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	location, ok := storage.DataStorage.Location.Get(id)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	writeObj(ctx, location)
}

func GetLocatioAvgnEndpoint(ctx *fasthttp.RequestCtx) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	params := ctx.QueryArgs()

	//  Parse fromDate parameter
	var fromDate int64
	if params.Has("fromDate") {
		if _, err := params.GetUint("fromDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("fromDate")
		fromDate = int64(tmp)
	}

	//  Parse toDate parameter
	var toDate int64
	if params.Has("toDate") {
		if _, err := params.GetUint("toDate"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toDate")
		toDate = int64(tmp)
	}

	//  Parse fromAge parameter
	var fromAge int64
	if params.Has("fromAge") {
		if _, err := params.GetUint("fromAge"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("fromAge")
		fromAge = int64(tmp)
	}

	//  Parse toAge parameter
	var toAge int64
	if params.Has("toAge") {
		if _, err := params.GetUint("toAge"); err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		tmp, _ := params.GetUint("toAge")
		toAge = int64(tmp)
	}

	//  Parse gender parameter
	var gender string
	if params.Has("gender") {
		gender = string(params.Peek("gender"))
		if gender != "m" && gender != "f" {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	}

	visits, ok := storage.DataStorage.VisitList.GetByLocation(id, storage.DataStorage)
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	var sum int32
	var count int32
	for _, visit := range visits {
		if visit.UserID == nil {
			continue
		}
		user, ok := storage.DataStorage.User.Get(*visit.UserID)
		if !ok {
			continue
		}
		if params.Has("fromDate") && fromDate >= (*visit.VisitedAt) {
			continue
		}
		if params.Has("toDate") && toDate <= (*visit.VisitedAt) {
			continue
		}

		if params.Has("fromAge") &&
			time.Now().AddDate(-int(fromAge), 0, 0).Unix() <= (*user.BirthDate) {
			continue
		}
		if params.Has("toAge") &&
			time.Now().AddDate(-int(toAge), 0, 0).Unix() >= (*user.BirthDate) {
			continue
		}
		if params.Has("gender") && gender != (*user.Gender) {
			continue
		}
		count++
		sum += int32(*visit.Mark)
	}

	var avg float64
	if count > 0 {
		avg = Round(float64(sum)/float64(count), 0.5, 5)
	}

	answ := fmt.Sprintf(`{"avg": %.5f}`, avg)
	writeStr(ctx, answ)
}

func UpdateLocationEndpoint(ctx *fasthttp.RequestCtx) {
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
	var location model.Location
	_ = ffjson.Unmarshal(bytes, &location)

	err = storage.DataStorage.Location.Update(id, &location)
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
	_ = ffjson.Unmarshal(ctx.PostBody(), &location)

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
