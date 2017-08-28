package handler

import (
	"encoding/json"
	"net/http"

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

		avg := 0
		for _, visit := range visits {
			// TODO filter visit
			if data.Gender != nil && (*data.Gender) == (*visit.User.Gender) {

			}
		}

		var resp schema.ResponceLocationVisits
		resp.Avg = float32(avg)
		json.NewEncoder(w).Encode(resp)
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
