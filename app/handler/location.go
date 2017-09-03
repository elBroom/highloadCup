package handler

import (
	"encoding/json"
	"net/http"

	"fmt"

	"strconv"

	"time"

	"math"
	"net/url"

	"io/ioutil"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/elBroom/highloadCup/app/workers"
)

func GetLocationEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {

			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		location, ok := storage.DataStorage.Location.Get(id)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		json.NewEncoder(w).Encode(location)
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func GetLocatioAvgnEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		params, err := url.ParseQuery(req.URL.RawQuery)
		//  Parse fromDate parameter
		fromDateStr := params.Get("fromDate")
		fromDate := int64(0)

		if fromDateStr != "" {
			if fromDate, err = strconv.ParseInt(fromDateStr, 10, 64); err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return nil
			}
		}

		//  Parse toDate parameter
		toDateStr := params.Get("toDate")
		toDate := int64(math.MaxInt64)

		if toDateStr != "" {
			if toDate, err = strconv.ParseInt(toDateStr, 10, 64); err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return nil
			}
		}

		//  Parse fromAge parameter
		fromAgeStr := params.Get("fromAge")
		fromAge := int64(0)

		if fromAgeStr != "" {
			if fromAge, err = strconv.ParseInt(fromAgeStr, 10, 64); err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return nil
			}
		}

		//  Parse toAge parameter
		toAgeStr := params.Get("toAge")
		toAge := int64(math.MaxInt64)

		if toAgeStr != "" {
			if toAge, err = strconv.ParseInt(toAgeStr, 10, 64); err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return nil
			}
		}

		//  Parse gender parameter
		gender := params.Get("gender")
		if gender != "" && gender != "m" && gender != "f" {
			http.Error(w, "", http.StatusBadRequest)
			return nil
		}

		visits, ok := storage.DataStorage.VisitList.GetByLocation(id, storage.DataStorage)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return nil
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
			if fromDateStr != "" && fromDate >= (*visit.VisitedAt) {
				continue
			}
			if toDateStr != "" && toDate <= (*visit.VisitedAt) {
				continue
			}

			if fromAgeStr != "" &&
				time.Now().AddDate(-int(fromAge), 0, 0).Unix() <= (*user.BirthDate) {
				continue
			}
			if toAgeStr != "" &&
				time.Now().AddDate(-int(toAge), 0, 0).Unix() >= (*user.BirthDate) {
				continue
			}
			if gender != "" && gender != (*user.Gender) {
				continue
			}
			count++
			sum += int32(*visit.Mark)
		}

		var avg float64
		if count > 0 {
			avg = Round(float64(sum)/float64(count), 0.5, 5)
		}

		answ := []byte(fmt.Sprintf(`{"avg": %.5f}`, avg))
		w.Write(answ)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", strconv.Itoa(len(answ)))
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func UpdateLocationEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		bytes, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		ok := CheckNull(bytes)
		if ok {
			http.Error(w, "", http.StatusBadRequest)
			return nil
		}
		var location model.Location
		_ = json.Unmarshal(bytes, &location)

		err = storage.DataStorage.Location.Update(id, &location)
		if err != nil {
			if err == storage.ErrDoesNotExist {
				http.Error(w, "", http.StatusNotFound)
			} else {
				http.Error(w, "", http.StatusBadRequest)
			}
			return nil
		}
		w.Write([]byte("{}"))
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func CreateLocationEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		var location model.Location
		defer req.Body.Close()
		_ = json.NewDecoder(req.Body).Decode(&location)

		err := storage.DataStorage.Location.Add(&location)
		if err != nil {

			if err == storage.ErrDoesNotExist {
				http.Error(w, "", http.StatusNotFound)
			} else {
				http.Error(w, "", http.StatusBadRequest)
			}
			return nil
		}

		w.Write([]byte("{}"))
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}
