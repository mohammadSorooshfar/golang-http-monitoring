package store

import (
	"context"

	"github.com/mohammadSorooshfar/golang-http-monitoring/model"
)

type AuthInMemory struct {
	Users map[string]model.Auth
}

func NewAuthInMemory() *AuthInMemory {
	return &AuthInMemory{
		Users: make(map[string]model.Auth),
	}
}

func (m *AuthInMemory) Save(_ context.Context, s model.Auth) error {
	if _, ok := m.Users[s.Username]; ok {
		return DuplicateUserError{
			Username: s.Username,
		}
	}

	m.Users[s.Username] = s

	return nil
}

func (m *AuthInMemory) Get(_ context.Context, user model.Auth) (model.Auth, error) {
	s, ok := m.Users[user.Username]
	if ok {
		if s.Password == user.Password {

			return s, nil
		}
	}

	return s, UserNotFoundError{
		Username: user.Username,
	}
}
