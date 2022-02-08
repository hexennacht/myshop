package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName       string
	Username       string `gorm:"index:idx_login_search,uniqueIndex,"`
	Email          string `gorm:"index:idx_login_search,uniqueIndex"`
	Password       string
	PhoneNumber    int32
	BirthDayDate   time.Time
	ProfilePicture string
	Status         int
}
