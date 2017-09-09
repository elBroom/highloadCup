package storage

import (
	"sync"

	"log"

	"github.com/elBroom/highloadCup/app"
	"github.com/elBroom/highloadCup/app/model"
)

type Location struct {
	mx       sync.RWMutex
	location [CountLocation]*model.Location
}

func (l *Location) Add(location *model.Location) error {
	if location.ID == nil || location.Place == nil || location.Country == nil || location.City == nil ||
		location.Distance == nil {
		return ErrRequiredFields
	}

	if *(location.ID) > CountLocation {
		log.Printf("Big index location: %d", *(location.ID))
		return nil
	}

	val := l.location[*(location.ID)]
	if val != nil {
		return ErrAlreadyExist
	}

	l.mx.Lock()
	defer l.mx.Unlock()
	l.location[*(location.ID)] = location

	if app.Phase == 2 {
		go DataStorage.VisitList.AddEmptyForLocation(*(location.ID))
	} else {
		DataStorage.VisitList.AddEmptyForLocation(*(location.ID))
	}
	return nil
}

func (l *Location) Update(id uint32, new_location *model.Location) error {
	if new_location.ID != nil {
		return ErrIDInUpdate
	}

	location := l.location[id]
	if location == nil {
		return ErrDoesNotExist
	}

	l.mx.Lock()
	defer l.mx.Unlock()
	if new_location.Place != nil {
		location.Place = new_location.Place
	}
	if new_location.Country != nil {
		location.Country = new_location.Country
	}
	if new_location.City != nil {
		location.City = new_location.City
	}
	if new_location.Distance != nil {
		location.Distance = new_location.Distance
	}
	return nil
}

func (l *Location) Get(id uint32) (*model.Location, bool) {
	//l.mx.RLock()
	//defer l.mx.RUnlock()

	location := l.location[id]
	return location, location != nil
}
