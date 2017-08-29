package handler

import (
	"net/http"

	"encoding/json"

	"io/ioutil"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
	"github.com/elBroom/highloadCup/app/workers"
)

func GetVisitEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {

			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		visit, ok := storage.DataStorage.Visit.Get(id)
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		json.NewEncoder(w).Encode(visit)
		return nil
	}, workers.TimeOut)

	checkTimeout(w, err)
}

func UpdateVisitEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		id, err := parseID(req)
		if err != nil {

			http.Error(w, "", http.StatusNotFound)
			return nil
		}

		var visit model.Visit
		defer req.Body.Close()
		_ = json.NewDecoder(req.Body).Decode(&visit)

		bytes, _ := ioutil.ReadAll(req.Body)
		ok := CheckNull(bytes)
		if ok {
			http.Error(w, "", http.StatusBadRequest)
			return nil
		}

		err = storage.DataStorage.Visit.Update(id, &visit, storage.DataStorage)
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

func CreateVisitEndpoint(w http.ResponseWriter, req *http.Request) {
	_, err := workers.Wp.AddTaskSyncTimed(func() interface{} {
		var visit model.Visit
		defer req.Body.Close()
		_ = json.NewDecoder(req.Body).Decode(&visit)

		err := storage.DataStorage.Visit.Add(&visit, storage.DataStorage)
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
