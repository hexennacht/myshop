package routes

import (
	"github.com/hexennacht/myshop/user/handler/middleware"
	"github.com/labstack/echo/v4"
	middlewareEcho "github.com/labstack/echo/v4/middleware"

	"github.com/hexennacht/myshop/user/handler"
)

type Routes interface {
	RegisterMiddleware() Routes
	PublicRoute() Routes
	AuthorizedRoute() Routes
}

type route struct {
	jwtKeySecret string
	echo         *echo.Echo
	userHandler  handler.Handler
}

func NewRoutes(jwtKeySecret string, echo *echo.Echo, userHandler handler.Handler) Routes {
	return &route{jwtKeySecret: jwtKeySecret, echo: echo, userHandler: userHandler}
}

func (r *route) RegisterMiddleware() Routes {
	r.echo.Use(middlewareEcho.Logger())
	return r
}

func (r *route) PublicRoute() Routes {
	group := r.echo.Group("/auth")
	group.POST("/login", r.userHandler.UserLogin)
	group.POST("/register", r.userHandler.UserRegister)

	return r
}

func (r *route) AuthorizedRoute() Routes {
	group := r.echo.Group("/user", middleware.ValidateToken())
	group.GET("/", r.userHandler.GetUserDetail, middleware.Contexter)
	group.POST("/", r.userHandler.UpdatePassword, middleware.Contexter)
	return r
}
