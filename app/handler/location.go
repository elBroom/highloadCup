package handler

import (
	"encoding/json"
	"net/http"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/golang/glog"
)

func GetLocationEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	location, ok := storage.DataStorage.Location.Get(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(location)
}

func GetLocatioAvgnEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	var data schema.RequestLocationVisits
	_ = json.NewDecoder(req.Body).Decode(&data)
	defer req.Body.Close()

	visits, ok := storage.DataStorage.VisitList.GetByLocation(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
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
}

func UpdateLocationEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
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
		return
	}
	w.Write([]byte("{}"))
}

func CreateLocationEndpoint(w http.ResponseWriter, req *http.Request) {
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
		return
	}

	w.Write([]byte("{}"))
}
