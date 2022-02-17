package handler

import (
	"context"
	"github.com/hexennacht/myshop/user/helper"
	"github.com/hexennacht/myshop/user/module/entity"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"

	"github.com/hexennacht/myshop/user/module/user"
)

type Handler interface {
	UserLogin(c echo.Context) error
	UserRegister(c echo.Context) error
	GetUserDetail(c echo.Context) error
	UpdatePassword(c echo.Context) error
}

type handler struct {
	mod user.Module
}

func NewUserHandler(userModule user.Module) Handler {
	return &handler{mod: userModule}
}

func (h *handler) UserLogin(c echo.Context) error {
	var req entity.UserLogin

	if err := c.Bind(&req); err != nil {
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "request malformed",
				"error":   err,
				"code":    400,
			})
		}
	}

	resp, err := h.mod.UserLogin(context.Background(), &req)
	if err != nil {
		var errorMessages []string
		for _, err := range err {
			errorMessages = append(errorMessages, strings.Replace(err.Error(), "UserLogin.", "", -1))
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "request malformed",
			"error":   errorMessages,
			"code":    400,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user logged in",
		"data":    resp,
		"code":    200,
	})
}

func (h *handler) UserRegister(c echo.Context) error {
	var req entity.CreateUser
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "request malformed",
			"error":   err,
			"code":    400,
		})
	}
	if errs := h.mod.RegisterUser(context.Background(), &req); errs != nil {
		var errorMessages []string
		for _, err := range errs {
			errorMessages = append(errorMessages, strings.Replace(err.Error(), "CreateUser.", "", -1))
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "request malformed",
			"error":   errorMessages,
			"code":    400,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user success created",
		"code":    200,
	})
}

func (h *handler) GetUserDetail(c echo.Context) error {
	ctx, ok := helper.FromEchoContext(c)
	if !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden access",
			"code":    403,
		})
	}

	userCtx, ok := helper.ParseContext(ctx)
	if !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden access",
			"code":    403,
		})
	}
	user, err := h.mod.GetUser(ctx, &entity.UserLogin{User: userCtx.Username})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "internal server error",
			"error":   err.Error(),
			"code":    500,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    user,
		"code":    200,
	})
}

func (h *handler) UpdatePassword(c echo.Context) error {
	ctx, ok := helper.FromEchoContext(c)
	if !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden access",
			"code":    403,
		})
	}

	userCtx, ok := helper.ParseContext(ctx)
	if !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "forbidden access",
			"code":    403,
		})
	}

	var req entity.UserUpdatePasswordRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "request malformed",
			"error":   err,
			"code":    400,
		})
	}

	req.Username = userCtx.Username
	if errs := h.mod.UpdateUserPassword(context.Background(), &req); errs != nil {
		var errorMessages []string
		for _, err := range errs {
			errorMessages = append(errorMessages, strings.Replace(err.Error(), "UserUpdatePasswordRequest.", "", -1))
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "request malformed",
			"error":   errorMessages,
			"code":    400,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user update password success",
		"code":    200,
	})
}
