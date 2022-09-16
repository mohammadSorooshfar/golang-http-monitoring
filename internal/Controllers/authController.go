package controllers

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var ErrInvalidSigningMethod = errors.New("unexpected signing method")

const TokenHeaderLen = 2

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenHeader := strings.Fields(c.Request().Header.Get("Authorization"))

		if len(tokenHeader) != TokenHeaderLen {
			return echo.ErrUnauthorized
		}

		tokenString := tokenHeader[1]

		token, err := jwt.ParseWithClaims(tokenString, new(jwt.RegisteredClaims),
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, ErrInvalidSigningMethod
				}

				return []byte("secret"), nil
			})
		if err != nil {
			return echo.ErrUnauthorized
		}

		if !token.Valid {
			return echo.ErrUnauthorized
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {

			return echo.ErrUnauthorized
		}

		if err := claims.Valid(); err != nil {

			return echo.ErrUnauthorized
		}
		c.Set("name", claims.Subject)
		c.Set("id", claims.ID)

		return next(c)
	}
}
