package storage

import (
	"sync"

	"log"

	"math"
	"time"

	"github.com/elBroom/highloadCup/app/model"
	"github.com/elBroom/highloadCup/app/schema"
	"github.com/google/btree"
)

var deqree = 4

type VisitList struct {
	mxU      sync.RWMutex
	mxL      sync.RWMutex
	user     [CountUser](*btree.BTree)
	location [CountLocation](*btree.BTree)
}

type VisitItem struct {
	*model.Visit
}

func (v VisitItem) Less(i btree.Item) bool {
	return *v.VisitedAt < *i.(VisitItem).VisitedAt
}

func (vl *VisitList) Add(visit *model.Visit) {
	if *(visit.LocationID) > CountLocation {
		log.Printf("Big index location from visit: %d", *(visit.ID))
	}

	if *(visit.UserID) > CountUser {
		log.Printf("Big index user from visit: %d", *(visit.ID))
	}

	vl.AddLocation(visit)
	vl.AddUser(visit)
}

func (vl *VisitList) AddLocation(visit *model.Visit) {
	vl.mxL.Lock()
	if *(visit.LocationID) <= CountLocation {
		if vl.location[*(visit.LocationID)] == nil {
			vl.location[*(visit.LocationID)] = btree.New(deqree)
		}
		vl.location[*(visit.LocationID)].ReplaceOrInsert(VisitItem{visit})
	}
	vl.mxL.Unlock()
}

func (vl *VisitList) AddUser(visit *model.Visit) {
	vl.mxU.Lock()
	if *(visit.UserID) <= CountUser {
		if vl.user[*(visit.UserID)] == nil {
			vl.user[*(visit.UserID)] = btree.New(deqree)
		}
		vl.user[*(visit.UserID)].ReplaceOrInsert(VisitItem{visit})
	}
	vl.mxU.Unlock()
}

func (vl *VisitList) AddEmptyForLocation(id uint32) {
	vl.mxL.Lock()
	defer vl.mxL.Unlock()

	visits := vl.location[id]
	if visits == nil {
		vl.location[id] = btree.New(deqree)
	}
}

func (vl *VisitList) AddEmptyForUser(id uint32) {
	vl.mxU.Lock()
	defer vl.mxU.Unlock()

	visits := vl.user[id]
	if visits == nil {
		vl.user[id] = btree.New(deqree)
	}
}

func (vl *VisitList) Delete(visit *model.Visit, new_visit *model.Visit) {
	if new_visit.LocationID != nil && *visit.LocationID != *new_visit.LocationID ||
		new_visit.VisitedAt != nil && *visit.VisitedAt != *new_visit.VisitedAt {
		vl.mxL.Lock()
		// Delete visit from old location
		vl.location[*(visit.LocationID)].Delete(VisitItem{visit})
		vl.mxL.Unlock()
	}

	if new_visit.UserID != nil && *visit.UserID != *new_visit.UserID ||
		new_visit.VisitedAt != nil && *visit.VisitedAt != *new_visit.VisitedAt {
		vl.mxU.Lock()
		// Delete visit from old user
		vl.user[*(visit.UserID)].Delete(VisitItem{visit})
		vl.mxU.Unlock()
	}
}

func (vl *VisitList) GetByLocation(id uint32,
	fromDate *int64,
	toDate *int64,
	fromAge *int64,
	toAge *int64,
	gender *string) *float64 {
	//vl.mx.RLock()
	//defer vl.mx.RUnlock()

	if vl.location[id] == nil {
		return nil
	}

	now := time.Now()
	var sum int32
	var count int32
	iterator := func(i btree.Item) bool {
		user, ok := DataStorage.User.Get(*i.(VisitItem).UserID)
		if !ok {
			return true
		}
		if fromAge != nil && now.AddDate(-int(*fromAge), 0, 0).Unix() <= (*user.BirthDate) {
			return true
		}
		if toAge != nil && now.AddDate(-int(*toAge), 0, 0).Unix() >= (*user.BirthDate) {
			return true
		}
		if gender != nil && *gender != (*user.Gender) {
			return true
		}
		count++
		sum += int32(*i.(VisitItem).Mark)
		return true
	}

	if fromDate == nil && toDate == nil {
		vl.location[id].Ascend(iterator)
	} else if fromDate == nil && toDate != nil {
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vl.location[id].AscendLessThan(lessThan, iterator)
	} else if fromDate != nil && toDate == nil {
		(*fromDate)++
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: fromDate}}
		vl.location[id].AscendGreaterOrEqual(greaterOrEqual, iterator)
	} else {
		(*fromDate)++
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: fromDate}}
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vl.location[id].AscendRange(greaterOrEqual, lessThan, iterator)
	}

	var avg float64
	if count > 0 {
		avg = Round(float64(sum)/float64(count), 0.5, 5)
	}
	return &avg
}

func (vl *VisitList) GetByUser(id uint32,
	fromDate *int64,
	toDate *int64,
	country *string,
	toDistance *uint32,
	visits *[]*schema.ResponceUserVisit) bool {
	//vl.mx.RLock()
	//defer vl.mx.RUnlock()

	if vl.user[id] == nil {
		return false
	}

	iterator := func(i btree.Item) bool {
		location, ok := DataStorage.Location.Get(*i.(VisitItem).LocationID)
		if !ok {
			return true
		}
		if country != nil && *country != (*location.Country) {
			return true
		}
		if toDistance != nil && *toDistance <= (*location.Distance) {
			return true
		}
		var item schema.ResponceUserVisit
		item.Mark = i.(VisitItem).Mark
		item.Visited_at = i.(VisitItem).VisitedAt
		item.Place = location.Place
		*visits = append(*visits, &item)
		return true
	}

	if fromDate == nil && toDate == nil {
		vl.user[id].Ascend(iterator)
	} else if fromDate == nil && toDate != nil {
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vl.user[id].AscendLessThan(lessThan, iterator)
	} else if fromDate != nil && toDate == nil {
		(*fromDate)++
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: fromDate}}
		vl.user[id].AscendGreaterOrEqual(greaterOrEqual, iterator)
	} else {
		(*fromDate)++
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: fromDate}}
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vl.user[id].AscendRange(greaterOrEqual, lessThan, iterator)
	}

	return true
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
