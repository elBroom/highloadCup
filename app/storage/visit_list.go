package storage

import (
	"sync"

	"github.com/elBroom/highloadCup/app/model"
)

type VisitList struct {
	mx       sync.RWMutex
	user     map[uint32]([]*model.Visit)
	location map[uint32]([]*model.Visit)
}

func (vl *VisitList) Add(visit *model.Visit) error {
	vl.mx.Lock()
	defer vl.mx.Unlock()
	location_visit_list, ok := vl.location[*(visit.LocationID)]
	if ok {
		vl.location[*(visit.LocationID)] = append(location_visit_list, visit)
	} else {
		vl.location[*(visit.LocationID)] = []*model.Visit{visit}
	}

	user_visit_list, ok := vl.user[*(visit.UserID)]
	if ok {
		vl.user[*(visit.UserID)] = append(user_visit_list, visit)
	} else {
		vl.user[*(visit.UserID)] = []*model.Visit{visit}
	}
	return nil
}

func (vl *VisitList) Update(old_visit *model.Visit, new_visit *model.Visit) error {
	vl.mx.Lock()
	defer vl.mx.Unlock()

	if old_visit.LocationID != new_visit.LocationID {
		// Delete old location
		location_old, ok := vl.location[*(old_visit.LocationID)]
		if ok {
			for index, visit := range location_old {
				if visit.ID == old_visit.ID {
					location_old = append(location_old[:index], location_old[index+1:]...)
					break
				}
			}
		}
		// Add new location
		location_new, ok := vl.location[*(new_visit.LocationID)]
		if ok {
			vl.location[*(new_visit.LocationID)] = append(location_new, new_visit)
		} else {
			vl.location[*(new_visit.LocationID)] = []*model.Visit{new_visit}
		}
	}

	if old_visit.UserID != new_visit.UserID {
		// Delete old user
		user_old, ok := vl.location[*(old_visit.LocationID)]
		if ok {
			for index, visit := range user_old {
				if visit.ID == old_visit.ID {
					user_old = append(user_old[:index], user_old[index+1:]...)
					break
				}
			}
		}
		// Add old user
		user_new, ok := vl.location[*(new_visit.LocationID)]
		if ok {
			vl.location[*(new_visit.LocationID)] = append(user_new, new_visit)
		} else {
			vl.location[*(new_visit.LocationID)] = []*model.Visit{new_visit}
		}
	}
	return nil
}

func (vl *VisitList) GetByLocation(id uint32) ([]*model.Visit, bool) {
	vl.mx.RLock()
	defer vl.mx.RUnlock()

	visits, ok := vl.location[id]

	if ok {
		return visits, ok
	}
	return nil, ok
}

func (vl *VisitList) GetByUser(id uint32) ([]*model.Visit, bool) {
	vl.mx.RLock()
	defer vl.mx.RUnlock()

	visits, ok := vl.user[id]

	if ok {
		return visits, ok
	}
	return nil, ok
}
