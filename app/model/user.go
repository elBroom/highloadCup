package model

import (
	"errors"
)

var ErrNullField = errors.New("Field value is null")

type User struct {
	ID        *uint32 `json:"id"`
	Email     *string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Gender    *string `json:"gender"`
	BirthDate *int32  `json:"birth_date"`
}

type Users struct {
	Users []User `json:"users"`
}
