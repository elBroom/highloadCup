package handler

import (
	"net/http"

	"encoding/json"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/golang/glog"
)

func GetVisitEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	visit, ok := storage.DataStorage.Visit.Get(id)
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(visit)
}

func UpdateVisitEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	id, err := parseID(req)
	if err != nil {
		glog.Infoln(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	var visit model.Visit
	_ = json.NewDecoder(req.Body).Decode(&visit)
	defer req.Body.Close()

	err = storage.DataStorage.Visit.Update(id, &visit, storage.DataStorage)
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

func CreateVisitEndpoint(w http.ResponseWriter, req *http.Request) {
	glog.Infoln(req.Method, req.RequestURI)

	var visit model.Visit
	_ = json.NewDecoder(req.Body).Decode(&visit)
	defer req.Body.Close()

	err := storage.DataStorage.Visit.Add(&visit, storage.DataStorage)
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
