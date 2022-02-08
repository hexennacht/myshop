package repository

import (
	"github.com/hexennacht/myshop/user/module"
	"github.com/hexennacht/myshop/user/module/entity"
	"gorm.io/gorm"
	"time"
)

func TransformToUser(req *entity.CreateUser) *User {
	return &User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		FullName:       req.FullName,
		Username:       req.Username,
		Email:          req.Email,
		Password:       req.HassedPassword,
		PhoneNumber:    req.PhoneNumber,
		BirthDayDate:   req.BirthDayDate,
		ProfilePicture: req.ProfilePicture,
		Status:         int(module.UserStatusActive),
	}
}

func (u *User) TransformToEntity() *entity.User {
	return &entity.User{
		FullName:       u.FullName,
		Username:       u.Username,
		Email:          u.Email,
		PhoneNumber:    u.PhoneNumber,
		Password:       u.Password,
		BirthDayDate:   u.BirthDayDate,
		ProfilePicture: u.ProfilePicture,
		Status:         module.UserStatus(u.Status),
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
