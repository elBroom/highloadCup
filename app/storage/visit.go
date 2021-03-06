package storage

import (
	"sync"

	"log"

	"github.com/elBroom/highloadCup/app"
	"github.com/elBroom/highloadCup/app/model"
)

type Visit struct {
	mx    sync.RWMutex
	visit [CountVisit]*model.Visit
}

func (v *Visit) Add(visit *model.Visit) error {
	if visit.ID == nil || visit.LocationID == nil || visit.UserID == nil ||
		visit.VisitedAt == nil || visit.Mark == nil || (*visit.Mark) < 0 {
		return ErrRequiredFields
	}

	if *(visit.ID) > CountVisit {
		log.Printf("Big index visit: %d", *(visit.ID))
		return nil
	}

	val := v.visit[*(visit.ID)]
	if val != nil {
		return ErrAlreadyExist
	}

	v.mx.Lock()
	defer v.mx.Unlock()
	v.visit[*(visit.ID)] = visit

	if app.Phase == 2 {
		go DataStorage.VisitList.Add(visit)
	} else {
		DataStorage.VisitList.Add(visit)
	}
	return nil
}

func (v *Visit) Update(id uint32, new_visit *model.Visit) error {
	if new_visit.ID != nil {
		return ErrIDInUpdate
	}

	visit := v.visit[id]
	if visit == nil {
		return ErrDoesNotExist
	}

	v.mx.Lock()
	defer v.mx.Unlock()

	isChangeLocation := new_visit.LocationID != nil && (*visit.LocationID) != (*new_visit.LocationID)
	isChangeUser := new_visit.UserID != nil && (*visit.UserID) != (*new_visit.UserID)
	isChangeVisitAt := new_visit.VisitedAt != nil && (*visit.VisitedAt) != (*new_visit.VisitedAt)
	if isChangeLocation || isChangeUser || isChangeVisitAt {
		DataStorage.VisitList.Delete(visit, new_visit)
	}
	if isChangeLocation {
		visit.LocationID = new_visit.LocationID
	}
	if isChangeUser {
		visit.UserID = new_visit.UserID
	}
	if isChangeVisitAt {
		visit.VisitedAt = new_visit.VisitedAt
	}
	if new_visit.Mark != nil && (*new_visit.Mark) >= 0 {
		visit.Mark = new_visit.Mark
	}

	go func() {
		if isChangeLocation || isChangeVisitAt {
			DataStorage.VisitList.AddLocation(visit)
		}
		if isChangeUser || isChangeVisitAt {
			DataStorage.VisitList.AddUser(visit)
		}
	}()

	return nil
}

func (v *Visit) Get(id uint32) (*model.Visit, bool) {
	//v.mx.RLock()
	//defer v.mx.RUnlock()

	visit := v.visit[id]
	return visit, visit != nil
}
