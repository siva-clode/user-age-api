package models

import "time"

type User struct {
	ID   int       `json:"id"`
	Name string    `json:"name" validate:"required"`
	DOB  time.Time `json:"dob" validate:"required"`
}
