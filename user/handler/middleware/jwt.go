package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/hexennacht/myshop/user/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.
type Claims struct {
	FullName string `json:"full_name"`
	UserCode string `json:"user_code"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func ValidateToken() echo.MiddlewareFunc {
	secret := []byte(config.Read().SecretJWT)
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    secret,
		SigningMethod: jwt.SigningMethodHS256.Name,
		Claims:        &Claims{},
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) { // copy paste from echo framework example
			token, err := jwt.ParseWithClaims(auth, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				method, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("Signing method invalid")
				}

				HMAC256 := jwt.SigningMethodHS256.Name
				if method.Name != HMAC256 {
					return nil, fmt.Errorf("Signing method invalid")
				}

				return []byte(config.Read().SecretJWT), nil
			})
			if err != nil {
				return nil, err
			}

			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			claims, ok := token.Claims.(*Claims)
			if ok {
				ctxUserCode := context.WithValue(context.Background(), "Username", claims.UserCode)
				ctxName := context.WithValue(ctxUserCode, "FullName", claims.FullName)
				ctx := context.WithValue(ctxName, "Role", claims.Role)

				c.Request().WithContext(ctx)
			}

			return token, nil
		},
	})
}
