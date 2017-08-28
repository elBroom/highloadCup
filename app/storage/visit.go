package storage

import (
	"sync"

	"github.com/elBroom/highloadCup/app/model"
)

type Visit struct {
	mx    sync.RWMutex
	visit map[uint32]*model.Visit
}

func (v *Visit) Add(visit *model.Visit, st *Storage) error {
	if visit.ID == nil || visit.LocationID == nil || visit.UserID == nil ||
		visit.VisitedAt == nil || visit.Mark == nil {
		return ErrRequiredFields
	}
	v.mx.Lock()
	defer v.mx.Unlock()
	_, ok := v.visit[*(visit.ID)]
	if ok {
		return ErrAlreadyExist
	}

	location, ok := st.Location.Get(*(visit.LocationID))
	if !ok {
		return ErrIDInUpdate
	}
	visit.Location = location

	user, ok := st.User.Get(*(visit.UserID))
	if !ok {
		return ErrIDInUpdate
	}
	visit.User = user

	v.visit[*(visit.ID)] = visit

	st.VisitList.Add(visit)
	return nil
}

func (v *Visit) Update(id uint32, new_visit *model.Visit, st *Storage) error {
	if new_visit.ID != nil {
		return ErrIDInUpdate
	}

	v.mx.Lock()
	defer v.mx.Unlock()
	old_visit, ok := v.visit[id]
	if !ok {
		return ErrDoesNotExist
	}

	var location *model.Location
	var user *model.User
	isChangeLocation := new_visit.LocationID != nil && old_visit.LocationID != new_visit.LocationID
	isChangeUser := new_visit.UserID != nil && old_visit.UserID != new_visit.UserID
	if isChangeLocation || isChangeUser {
		if isChangeLocation {
			location, ok = st.Location.Get(*(new_visit.LocationID))
			if !ok {
				return ErrIDInUpdate
			}
		}

		if isChangeUser {
			user, ok = st.User.Get(*(new_visit.UserID))
			if !ok {
				return ErrIDInUpdate
			}
		}

		st.VisitList.Update(old_visit, new_visit)
	}
	if isChangeLocation {
		old_visit.Location = location
		old_visit.LocationID = new_visit.LocationID
	}
	if isChangeUser {
		old_visit.User = user
		old_visit.UserID = new_visit.UserID
	}
	if new_visit.VisitedAt != nil {
		old_visit.VisitedAt = new_visit.VisitedAt
	}
	if new_visit.Mark != nil {
		old_visit.Mark = new_visit.Mark
	}
	return nil
}

func (v *Visit) Get(id uint32) (*model.Visit, bool) {
	v.mx.RLock()
	defer v.mx.RUnlock()

	visit, ok := v.visit[id]

	if ok {
		visit_ := *visit
		return &visit_, ok
	}
	return nil, ok
}
