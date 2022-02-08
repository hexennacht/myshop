package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hexennacht/myshop/user/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func Contexter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", -1)

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Signing method invalid")
			}

			HMAC512 := jwt.SigningMethodHS256.Name
			if method.Name != HMAC512 {
				return nil, fmt.Errorf("Signing method invalid")
			}

			return []byte(config.Read().SecretJWT), nil
		})
		if err != nil {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "forbidden",
				"data":    err.Error(),
				"code":    403,
			})
		}

		if token == nil {
			return next(c)
		}

		claims, ok := token.Claims.(*Claims)
		if ok {
			ctxUserCode := context.WithValue(context.Background(), "Username", claims.UserCode)
			ctxName := context.WithValue(ctxUserCode, "FullName", claims.FullName)
			ctx := context.WithValue(ctxName, "Role", claims.Role)

			c.Request().WithContext(ctx)
		}
		return next(c)
	}
}
