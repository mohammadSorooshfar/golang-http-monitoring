package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mohammadSorooshfar/golang-http-monitoring/model"
	"github.com/mohammadSorooshfar/golang-http-monitoring/request"
	"github.com/mohammadSorooshfar/golang-http-monitoring/store"
)

var ErrInvalidSigningMethod = errors.New("unexpected signing method")

const TokenHeaderLen = 2

type Auth struct {
	Store store.Auth
}

func (a Auth) Auth(next echo.HandlerFunc) echo.HandlerFunc {
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

		return next(c)
	}
}

func (a Auth) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.Login

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {

		return echo.ErrBadRequest
	}

	data := model.Auth{
		Username: req.Username,
		Password: req.Password,
	}
	err := a.Store.Login(ctx, data)
	if err != nil {
		var errNotFound store.UserNotFoundError
		if ok := errors.As(err, &errNotFound); ok {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError

	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "students-summer-2022",
		Subject:   data.Username,
		Audience:  []string{"admin"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        data.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tokenString)
}

func (a Auth) Signup(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.Login

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {

		return echo.ErrBadRequest
	}

	data := model.Auth{
		Username: req.Username,
		Password: req.Password,
	}
	err := a.Store.Singup(ctx, data)
	if err != nil {
		var errDuplicate store.DuplicateUserError
		if ok := errors.As(err, &errDuplicate); ok {
			return echo.ErrBadRequest
		}

		return echo.ErrInternalServerError

	}

	return c.JSON(http.StatusOK, data)
}

func (a Auth) Register(g *echo.Group) {
	g.POST("/login", a.Login)
	g.POST("/signup", a.Signup)
}
