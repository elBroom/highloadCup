package importer

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/storage"
)

var errUnsupportedFile = errors.New("Imported file is not supported")

func ImportDataFromZip() error {
	r, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		name := f.Name

		// cut off folder part
		i := strings.LastIndex(name, "/")
		if i != -1 {
			name = name[i:]
		}
		// cut off extension part
		i = strings.LastIndex(name, ".")
		if i != -1 {
			name = name[:i]
		}

		parts := strings.Split(name, "_")
		if name == "options" {
			continue
		}
		if len(parts) != 2 {
			return errUnsupportedFile
		}

		// add concurrent processing
		rc, err := f.Open()
		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}
		err = rc.Close()
		if err != nil {
			return err
		}

		store := storage.DataStorage
		switch parts[0] {
		case "users":
			err = importUsers(bytes, store)
		case "locations":
			err = importLocations(bytes, store)
		case "visits":
			err = importVisits(bytes, store)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func importUsers(b []byte, store *storage.Storage) error {
	var users model.Users
	err := json.Unmarshal(b, &users)

	if err != nil {
		return err
	}

	for _, user := range users.Users {
		user_ := user
		store.User.Add(&user_)
	}
	return nil
}

func importLocations(b []byte, store *storage.Storage) error {
	var locations model.Locations
	err := json.Unmarshal(b, &locations)

	if err != nil {
		return err
	}

	for _, location := range locations.Locations {
		location_ := location
		store.Location.Add(&location_)
	}
	return nil
}

func importVisits(b []byte, store *storage.Storage) error {
	var visits model.Visits
	err := json.Unmarshal(b, &visits)

	if err != nil {
		return err
	}

	for _, visit := range visits.Visits {
		visit_ := visit
		store.Visit.Add(&visit_, store)
	}
	return nil
}
