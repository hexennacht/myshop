package entity

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hexennacht/myshop/user/handler/middleware"
	"github.com/hexennacht/myshop/user/module"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FullName       string            `json:"full_name"`
	Username       string            `json:"username,omitempty"`
	Email          string            `json:"email,omitempty"`
	Password       string            `json:"-"`
	PhoneNumber    int32             `json:"phoneNumber,omitempty"`
	BirthDayDate   time.Time         `json:"birthDayDate"`
	ProfilePicture string            `json:"profilePicture,omitempty"`
	Status         module.UserStatus `json:"status,omitempty"`
	CreatedAt      time.Time         `json:"-"`
	UpdatedAt      time.Time         `json:"-"`
	Role           string            `json:"-"`
}

type UserLogin struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *UserLogin) ValidateRequest() []error {
	v := validator.New()
	err := v.Struct(u)
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok && len(validationErrors) > 0 {
		return []error{errors.New("failed to get error")}
	}

	var errs []error
	for _, e := range validationErrors {
		errs = append(errs, fmt.Errorf(e.Error()))
	}

	return errs
}

func (u *UserLogin) GenerateToken(user *User, expirationTime time.Time, secret []byte) (*UserLoginResponse, []error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &middleware.Claims{
		FullName: user.FullName,
		UserCode: user.Username,
		Role:     "user",
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			Issuer:    user.Username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, []error{err}
	}

	return &UserLoginResponse{
		Token:      tokenString,
		ExpireDate: expirationTime,
	}, nil
}

type UserLoginResponse struct {
	Token      string    `json:"token"`
	ExpireDate time.Time `json:"expire_date"`
}

type CreateUser struct {
	FullName             string    `validate:"required" json:"full_name,omitempty"`
	Username             string    `validate:"required" json:"username,omitempty"`
	Email                string    `validate:"required,email" json:"email,omitempty"`
	Password             string    `validate:"required,eqfield=ConfirmationPassword,gte=8,alphanum" json:"password,omitempty"`
	ConfirmationPassword string    `validate:"eqfield=Password" json:"confirmation_password,omitempty"`
	PhoneNumber          int32     `validate:"required" json:"phone_number,omitempty"`
	BirthDayDate         time.Time `validate:"required" json:"birth_day_date"`
	ProfilePicture       string    `validate:"url" json:"profile_picture,omitempty"`
	HassedPassword       string    `json:"-"`
}

func (c *CreateUser) ValidateRequest() []error {
	v := validator.New()
	err := v.Struct(c)
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok && len(validationErrors) > 0 {
		return []error{errors.New("failed to get error")}
	}

	var errs []error
	for _, e := range validationErrors {
		errs = append(errs, fmt.Errorf(e.Error()))
	}

	return errs
}

func (c *CreateUser) CreatePassword() error {
	c.HassedPassword = base64.StdEncoding.EncodeToString([]byte(c.Password))

	password, err := bcrypt.GenerateFromPassword([]byte(c.HassedPassword), bcrypt.DefaultCost)

	c.HassedPassword = string(password)

	return err
}

func (u *UserLogin) ComparePassword(pass string) error {
	password := base64.StdEncoding.EncodeToString([]byte(u.Password))

	return bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
}
