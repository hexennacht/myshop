package helper

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/hexennacht/myshop/user/handler/middleware"
	"github.com/labstack/echo/v4"
)

func FromEchoContext(c echo.Context) (context.Context, bool) {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, false
	}

	userData, ok := userToken.Claims.(*middleware.Claims)
	if !ok {
		return nil, false
	}

	baseContext := c.Request().Context()

	userDataContext := &UserContext{
		Username: userData.UserCode,
		FullName: userData.FullName,
	}

	return context.WithValue(baseContext, "UserContext", userDataContext), ok
}

type UserContext struct {
	Username string
	FullName string
}

func ParseContext(ctx context.Context) (*UserContext, bool) {
	userCtx, ok := ctx.Value("UserContext").(*UserContext)
	return userCtx, ok
}
