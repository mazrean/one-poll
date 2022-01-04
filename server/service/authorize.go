package service

import (
	"context"
	"errors"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
	ErrNoUser                = errors.New("no user")
)

type Authorization interface {
	Signup(ctx context.Context, name values.UserName, password values.UserPassword) (*domain.User, error)
	Login(ctx context.Context, name values.UserName, password values.UserPassword) (*domain.User, error)
	DeleteAccount(ctx context.Context, user *domain.User) error
}
