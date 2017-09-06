package storage

import (
	"sync"

	"github.com/elBroom/highloadCup/app/model"
)

type User struct {
	mx   sync.RWMutex
	user map[uint32]*model.User
}

func (u *User) Add(user *model.User) error {
	if user.BirthDate == nil || user.Email == nil || user.FirstName == nil ||
		user.LastName == nil || user.Gender == nil || user.ID == nil {
		return ErrRequiredFields
	}

	u.mx.Lock()
	defer u.mx.Unlock()
	_, ok := u.user[*(user.ID)]
	if ok {
		return ErrAlreadyExist
	}
	u.user[*(user.ID)] = user
	return nil
}

func (u *User) Update(id uint32, new_user *model.User) error {
	if new_user.ID != nil {
		return ErrIDInUpdate
	}

	u.mx.Lock()
	defer u.mx.Unlock()
	user, ok := u.user[id]
	if !ok {
		return ErrDoesNotExist
	}
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

// Return copy
func (u *User) Get(id uint32) (*model.User, bool) {
	//u.mx.RLock()
	//defer u.mx.RUnlock()

	user, ok := u.user[id]

	if ok {
		user_ := *user
		return &user_, ok
	}
	return nil, ok
}
