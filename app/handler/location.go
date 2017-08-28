package handler

import (
	"encoding/json"
	"net/http"

	"fmt"

	"strconv"

	"time"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/elBroom/highloadCup/app/workers"
	"github.com/golang/glog"
)

func GetLocationEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		glog.Infoln(req.Method, req.RequestURI)

		id, err := parseID(req)
		if err != nil {
			glog.Infoln(err)
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
		glog.Infoln(req.Method, req.RequestURI)

		id, err := parseID(req)
		if err != nil {
			glog.Infoln(err)
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		var data schema.RequestLocationVisits
		_ = json.NewDecoder(req.Body).Decode(&data)
		defer req.Body.Close()

		visits, ok := storage.DataStorage.VisitList.GetByLocation(id)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		var sum int32
		var count int32
		for _, visit := range visits {
			if data.FromDate != nil && (*data.FromDate) >= (*visit.VisitedAt) {
				continue
			}
			if data.ToDate != nil && (*data.ToDate) <= (*visit.VisitedAt) {
				continue
			}
			if data.FromAge != nil &&
				time.Now().AddDate(-(*data.FromAge), 0, 0).Unix() >= (*visit.User.BirthDate) {
				continue
			}
			if data.ToAge != nil &&
				time.Now().AddDate(-(*data.ToAge), 0, 0).Unix() <= (*visit.User.BirthDate) {
				continue
			}
			if data.Gender != nil && (*data.Gender) != (*visit.User.Gender) {
				continue
			}
			count++
			sum += int32(*visit.Mark)
		}

		var avg float64
		if count > 0 {
			avg = Round(float64(sum/count), 0.5, 5)
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
		glog.Infoln(req.Method, req.RequestURI)

		id, err := parseID(req)
		if err != nil {
			glog.Infoln(err)
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		var location model.Location
		_ = json.NewDecoder(req.Body).Decode(&location)
		defer req.Body.Close()

		err = storage.DataStorage.Location.Update(id, &location)
		if err != nil {
			glog.Infoln(err)
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
		glog.Infoln(req.Method, req.RequestURI)

		var location model.Location
		_ = json.NewDecoder(req.Body).Decode(&location)
		defer req.Body.Close()

		err := storage.DataStorage.Location.Add(&location)
		if err != nil {
			glog.Infoln(err)
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
