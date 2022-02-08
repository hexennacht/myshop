package user

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"time"

	_ "golang.org/x/sync/errgroup"

	"github.com/hexennacht/myshop/user/module/entity"
	userRepository "github.com/hexennacht/myshop/user/repository/user"
)

type Module interface {
	UserLogin(ctx context.Context, req *entity.UserLogin) (*entity.UserLoginResponse, []error)
	RegisterUser(ctx context.Context, req *entity.CreateUser) []error
	GetUser(ctx context.Context, req *entity.UserLogin) (*entity.User, error)
}

type module struct {
	repo          userRepository.Repository
	tokenLifeTime int64
	jwtSecret     string
}

func NewModule(repo userRepository.Repository, tokenLifeTime int64, jwtSecret string) Module {
	return &module{repo: repo, tokenLifeTime: tokenLifeTime, jwtSecret: jwtSecret}
}

func (m *module) GetUser(ctx context.Context, req *entity.UserLogin) (*entity.User, error) {
	user, err := m.repo.GetUser(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *module) UserLogin(ctx context.Context, req *entity.UserLogin) (*entity.UserLoginResponse, []error) {
	errs := req.ValidateRequest()
	if len(errs) > 0 {
		return nil, errs
	}
	user, err := m.repo.GetUser(ctx, req.User)
	if err != nil {
		return nil, []error{err}
	}

	if err := req.ComparePassword(user.Password); err != nil {
		return nil, []error{err}
	}

	return req.GenerateToken(user, time.Now().Add(time.Duration(m.tokenLifeTime)*time.Hour), []byte(m.jwtSecret))
}

func (m *module) RegisterUser(ctx context.Context, req *entity.CreateUser) []error {
	errs := req.ValidateRequest()
	if len(errs) > 0 {
		return errs
	}

	var emailConfirm, userNameConfirm *entity.User
	eg, newCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		user, err := m.repo.GetUser(newCtx, req.Username)
		if err != nil {
			return err
		}

		userNameConfirm = user

		return nil
	})

	eg.Go(func() error {
		user, err := m.repo.GetUser(newCtx, req.Email)
		if err != nil {
			return err
		}

		emailConfirm = user

		return nil
	})

	if err := eg.Wait(); err == nil {
		if userNameConfirm != nil {
			return []error{errors.New("username already registered")}
		}

		if emailConfirm != nil {
			return []error{errors.New("email already registered")}
		}
	}

	if err := req.CreatePassword(); err != nil {
		return []error{err}
	}

	if err := m.repo.CreateUser(ctx, req); errs != nil {
		return []error{err}
	}

	return nil
}
