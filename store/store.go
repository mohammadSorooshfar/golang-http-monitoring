package store

import (
	"context"
	"fmt"
	"github.com/mohammadSorooshfar/golang-http-monitoring/model"
)

type UserNotFoundError struct {
	Username string
}

func (err UserNotFoundError) Error() string {
	return fmt.Sprintf("User %s doesn't exist", err.Username)
}

type DuplicateUserError struct {
	Username string
}

func (err DuplicateUserError) Error() string {
	return fmt.Sprintf("User %s already exists", err.Username)
}

type Auth interface {
	Login(context.Context, model.Auth) error
	Singup(context.Context, model.Auth) error
}
