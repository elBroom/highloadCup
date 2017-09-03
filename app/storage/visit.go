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
		visit.VisitedAt == nil || visit.Mark == nil || (*visit.Mark) < 0 {
		return ErrRequiredFields
	}
	v.mx.Lock()
	defer v.mx.Unlock()
	_, ok := v.visit[*(visit.ID)]
	if ok {
		return ErrAlreadyExist
	}

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
	visit, ok := v.visit[id]
	if !ok {
		return ErrDoesNotExist
	}

	old_visit := *visit
	isChangeLocation := new_visit.LocationID != nil && (*visit.LocationID) != (*new_visit.LocationID)
	isChangeUser := new_visit.UserID != nil && (*visit.UserID) != (*new_visit.UserID)
	if isChangeLocation {
		visit.LocationID = new_visit.LocationID
	}
	if isChangeUser {
		visit.UserID = new_visit.UserID
	}
	if new_visit.VisitedAt != nil {
		visit.VisitedAt = new_visit.VisitedAt
	}
	if new_visit.Mark != nil && (*new_visit.Mark) >= 0 {
		visit.Mark = new_visit.Mark
	}
	if isChangeLocation || isChangeUser {
		st.VisitList.Update(&old_visit, visit)
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
