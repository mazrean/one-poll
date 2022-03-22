package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
	"github.com/mazrean/one-poll/service"
)

type Authorization struct {
	db             repository.DB
	userRepository repository.User
}

func NewAuthorization(db repository.DB, userRepository repository.User) *Authorization {
	return &Authorization{
		db:             db,
		userRepository: userRepository,
	}
}

func (a *Authorization) Signup(ctx context.Context, name values.UserName, password values.UserPassword) (*domain.User, error) {
	var user *domain.User
	err := a.db.Transaction(ctx, nil, func(ctx context.Context) error {
		_, err := a.userRepository.GetUserByName(ctx, name)
		if err == nil {
			return service.ErrUserAlreadyExists
		}
		if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("failed to get user: %w", err)
		}

		hashedPassword, err := password.Hash()
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user = domain.NewUser(
			values.NewUserID(),
			name,
			hashedPassword,
		)

		err = a.userRepository.CreateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return user, nil
}

func (a *Authorization) Login(ctx context.Context, name values.UserName, password values.UserPassword) (*domain.User, error) {
	user, err := a.userRepository.GetUserByName(ctx, name)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrInvalidUserOrPassword
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = user.GetHashedPassword().Compare(password)
	if errors.Is(err, values.ErrInvalidPassword) {
		return nil, service.ErrInvalidUserOrPassword
	}
	if err != nil {
		return nil, fmt.Errorf("failed to compare password: %w", err)
	}

	return user, nil
}

func (a *Authorization) DeleteAccount(ctx context.Context, user *domain.User) error {
	err := a.userRepository.DeleteUser(ctx, user.GetID())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
