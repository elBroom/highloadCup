package storage

import (
	"sync"

	"github.com/elBroom/highloadCup/app/model"
)

type Location struct {
	mx       sync.RWMutex
	location map[uint32]*model.Location
}

func (l *Location) Add(location *model.Location) error {
	if location.ID == nil || location.Place == nil || location.Country == nil || location.City == nil ||
		location.Distance == nil {
		return ErrRequiredFields
	}

	l.mx.Lock()
	defer l.mx.Unlock()
	_, ok := l.location[*(location.ID)]
	if ok {
		return ErrAlreadyExist
	}
	l.location[*(location.ID)] = location
	DataStorage.VisitList.AddEmptyForLocation(*(location.ID))
	return nil
}

func (l *Location) Update(id uint32, new_location *model.Location) error {
	if new_location.ID != nil {
		return ErrIDInUpdate
	}

	l.mx.Lock()
	defer l.mx.Unlock()
	location, ok := l.location[id]
	if !ok {
		return ErrDoesNotExist
	}
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

	location, ok := l.location[id]
	return location, ok
}
