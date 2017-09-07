package storage

import (
	"sync"

	"log"

	"github.com/elBroom/highloadCup/app/model"
)

type VisitList struct {
	mx       sync.RWMutex
	user     [CountUser]([]*model.Visit)
	location [CountLocation]([]*model.Visit)
}

func (vl *VisitList) Add(visit *model.Visit) error {
	if *(visit.LocationID) > CountLocation {
		log.Printf("Big index location from visit: %d", *(visit.ID))
	}

	if *(visit.UserID) > CountUser {
		log.Printf("Big index user from visit: %d", *(visit.ID))
	}

	vl.mx.Lock()
	defer vl.mx.Unlock()
	if *(visit.LocationID) <= CountLocation {
		location_visit_list := vl.location[*(visit.LocationID)]
		if location_visit_list != nil {
			vl.location[*(visit.LocationID)] = append(location_visit_list, visit)
		} else {
			vl.location[*(visit.LocationID)] = []*model.Visit{visit}
		}
	}

	if *(visit.UserID) <= CountUser {
		user_visit_list := vl.user[*(visit.UserID)]
		if user_visit_list != nil {
			vl.user[*(visit.UserID)] = append(user_visit_list, visit)
		} else {
			vl.user[*(visit.UserID)] = []*model.Visit{visit}
		}
	}
	return nil
}

func (vl *VisitList) AddEmptyForLocation(id uint32) {
	vl.mx.Lock()
	defer vl.mx.Unlock()

	visits := vl.location[id]
	if visits == nil {
		vl.location[id] = []*model.Visit{}
	}
}

func (vl *VisitList) AddEmptyForUser(id uint32) {
	vl.mx.Lock()
	defer vl.mx.Unlock()

	visits := vl.user[id]
	if visits == nil {
		vl.user[id] = []*model.Visit{}
	}
}

func (vl *VisitList) Update(old_visit *model.Visit, new_visit *model.Visit) error {
	vl.mx.Lock()
	defer vl.mx.Unlock()

	if new_visit.LocationID != nil && (*old_visit.LocationID) != (*new_visit.LocationID) {
		// Delete visit from old location
		visits := vl.location[*(old_visit.LocationID)]
		if visits != nil {
			for index, visit := range visits {
				if (*visit.ID) == (*old_visit.ID) {
					vl.location[*(old_visit.LocationID)] = append(visits[:index], visits[index+1:]...)
					break
				}
			}
		}
		// Add visit to new location
		visits = vl.location[*(new_visit.LocationID)]
		if visits != nil {
			vl.location[*(new_visit.LocationID)] = append(visits, new_visit)
		} else {
			vl.location[*(new_visit.LocationID)] = []*model.Visit{new_visit}
		}
	}

	if new_visit.UserID != nil && (*old_visit.UserID) != (*new_visit.UserID) {
		// Delete visit from old user
		old_visits := vl.user[*(old_visit.UserID)]
		if old_visits != nil {
			for index, visit := range old_visits {
				if (*visit.ID) == (*old_visit.ID) {
					vl.user[*(old_visit.UserID)] = append(old_visits[:index], old_visits[index+1:]...)
					break
				}
			}
		}
		// Add visit to new user
		new_visits := vl.user[*(new_visit.UserID)]
		if new_visits != nil {
			vl.user[*(new_visit.UserID)] = append(new_visits, new_visit)
		} else {
			vl.user[*(new_visit.UserID)] = []*model.Visit{new_visit}
		}
	}
	return nil
}

func (vl *VisitList) GetByLocation(id uint32) ([]*model.Visit, bool) {
	//vl.mx.RLock()
	//defer vl.mx.RUnlock()

	visits := vl.location[id]
	return visits, visits != nil
}

func (vl *VisitList) GetByUser(id uint32) ([]*model.Visit, bool) {
	//vl.mx.RLock()
	//defer vl.mx.RUnlock()

	visits := vl.user[id]
	return visits, visits != nil
}
