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

	if new_visit.LocationID != nil && (*old_visit.LocationID) != (*new_visit.LocationID) {
		// Delete visit from old location
		visits, ok := vl.location[*(old_visit.LocationID)]
		if ok {
			for index, visit := range visits {
				if (*visit.ID) == (*old_visit.ID) {
					vl.location[*(old_visit.LocationID)] = append(visits[:index], visits[index+1:]...)
					break
				}
			}
		}
		// Add visit to new location
		visits, ok = vl.location[*(new_visit.LocationID)]
		if ok {
			vl.location[*(new_visit.LocationID)] = append(visits, new_visit)
		} else {
			vl.location[*(new_visit.LocationID)] = []*model.Visit{new_visit}
		}
	}

	if new_visit.UserID != nil && (*old_visit.UserID) != (*new_visit.UserID) {
		// Delete visit from old user
		old_visits, ok := vl.user[*(old_visit.UserID)]
		if ok {
			for index, visit := range old_visits {
				if (*visit.ID) == (*old_visit.ID) {
					vl.user[*(old_visit.UserID)] = append(old_visits[:index], old_visits[index+1:]...)
					break
				}
			}
		}
		// Add visit to new user
		new_visits, ok := vl.user[*(new_visit.UserID)]
		if ok {
			vl.user[*(new_visit.UserID)] = append(new_visits, new_visit)
		} else {
			vl.user[*(new_visit.UserID)] = []*model.Visit{new_visit}
		}
	}
	return nil
}

func (vl *VisitList) GetByLocation(id uint32, st *Storage) ([]*model.Visit, bool) {
	//vl.mx.RLock()
	//defer vl.mx.RUnlock()

	visits, ok := vl.location[id]

	if ok {
		return visits, ok
	}

	_, ok = st.Location.Get(id)
	if ok {
		return []*model.Visit{}, ok
	}
	return nil, ok
}

func (vl *VisitList) GetByUser(id uint32, st *Storage) ([]*model.Visit, bool) {
	//vl.mx.RLock()
	//defer vl.mx.RUnlock()

	visits, ok := vl.user[id]
	if ok {
		return visits, ok
	}

	_, ok = st.User.Get(id)
	if ok {
		return []*model.Visit{}, ok
	}
	return nil, ok
}
