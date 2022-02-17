package userRepository

import (
	"context"

	"gorm.io/gorm"

	"github.com/hexennacht/myshop/user/module/entity"
	repo "github.com/hexennacht/myshop/user/repository"
)

type Repository interface {
	CreateUser(ctx context.Context, req *entity.CreateUser) error
	GetUser(ctx context.Context, userNameEmail string) (*entity.User, error)
	UpdateUserPassword(ctx context.Context, userName, password string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, req *entity.CreateUser) error {
	return r.db.WithContext(ctx).Model(new(repo.User)).Create(repo.TransformToUser(req)).Error
}

func (r *repository) GetUser(ctx context.Context, userNameEmail string) (*entity.User, error) {
	var user repo.User

	err := r.db.WithContext(ctx).Model(new(repo.User)).Where("username = ? OR email = ?", userNameEmail, userNameEmail).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user.TransformToEntity(), nil
}

func (r *repository) UpdateUserPassword(ctx context.Context, userName, password string) error {
	return r.db.WithContext(ctx).Model(new(repo.User)).Where("username", userName).Updates(map[string]interface{}{"password": password}).Error
}
