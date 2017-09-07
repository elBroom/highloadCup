package storage

import (
	"errors"

	"github.com/elBroom/highloadCup/app/model"
)

const CountUser = 1100000
const CountLocation = 1100000
const CountVisit = 11000000

var (
	ErrRequiredFields = errors.New("Not all required fields are filled")
	ErrAlreadyExist   = errors.New("Already exist")
	ErrDoesNotExist   = errors.New("Does not exist")
	ErrIDInUpdate     = errors.New("Update should not contain ID in the json object")
)

type Storage struct {
	User             *User
	Location         *Location
	Visit            *Visit
	VisitList        *VisitList
	CurrentTimestamp int32
}

func NewStorage() *Storage {
	return &Storage{
		User:     &User{user: [CountUser]*model.User{}},
		Location: &Location{location: [CountLocation]*model.Location{}},
		Visit:    &Visit{visit: [CountVisit]*model.Visit{}},
		VisitList: &VisitList{
			user:     [CountUser]([]*model.Visit){},
			location: [CountLocation]([]*model.Visit){},
		},
	}
}

var DataStorage = NewStorage()
