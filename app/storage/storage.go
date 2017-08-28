package storage

import (
	"errors"

	"github.com/elBroom/highloadCup/app/model"
)

const countUser = 1000074
const countLocation = 1000074
const countVisit = 10000460

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
		User:     &User{user: make(map[uint32]*model.User, countUser)},
		Location: &Location{location: make(map[uint32]*model.Location, countLocation)},
		Visit:    &Visit{visit: make(map[uint32]*model.Visit, countVisit)},
		VisitList: &VisitList{
			user:     make(map[uint32]([]*model.Visit), countUser),
			location: make(map[uint32]([]*model.Visit), countLocation),
		},
	}
}

var DataStorage = NewStorage()
