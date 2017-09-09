package storage

import (
	"sync"

	"log"

	"github.com/elBroom/highloadCup/app"
	"github.com/elBroom/highloadCup/app/model"
)

type User struct {
	mx   sync.RWMutex
	user [CountUser]*model.User
}

func (u *User) Add(user *model.User) error {
	if user.BirthDate == nil || user.Email == nil || user.FirstName == nil ||
		user.LastName == nil || user.Gender == nil || user.ID == nil {
		return ErrRequiredFields
	}

	if *(user.ID) > CountUser {
		log.Printf("Big index user: %d", *(user.ID))
		return nil
	}

	if u.user[*(user.ID)] != nil {
		return ErrAlreadyExist
	}

	u.mx.Lock()
	defer u.mx.Unlock()
	u.user[*(user.ID)] = user

	if app.Phase == 2 {
		go DataStorage.VisitList.AddEmptyForUser(*(user.ID))
	} else {
		DataStorage.VisitList.AddEmptyForUser(*(user.ID))
	}
	return nil
}

func (u *User) Update(id uint32, new_user *model.User) error {
	if new_user.ID != nil {
		return ErrIDInUpdate
	}

	user := u.user[id]
	if user == nil {
		return ErrDoesNotExist
	}

	u.mx.Lock()
	defer u.mx.Unlock()
	if new_user.BirthDate != nil {
		user.BirthDate = new_user.BirthDate
	}
	if new_user.Email != nil {
		user.Email = new_user.Email
	}
	if new_user.FirstName != nil {
		user.FirstName = new_user.FirstName
	}
	if new_user.LastName != nil {
		user.LastName = new_user.LastName
	}
	if new_user.Gender != nil {
		user.Gender = new_user.Gender
	}
	return nil
}

func (u *User) Get(id uint32) (*model.User, bool) {
	//u.mx.RLock()
	//defer u.mx.RUnlock()

	user := u.user[id]
	return user, user != nil
}
