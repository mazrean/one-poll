package values

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserID             uuid.UUID
	UserName           string
	UserPassword       []byte
	UserHashedPassword []byte
)

func NewUserID() UserID {
	return UserID(uuid.New())
}

func NewUserIDFromUUID(id uuid.UUID) UserID {
	return UserID(id)
}

func NewUserName(name string) UserName {
	return UserName(name)
}

func NewUserPassword(password []byte) UserPassword {
	return UserPassword(password)
}

func (up UserPassword) Hash() (UserHashedPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(up), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return UserHashedPassword(hashedPassword), nil
}

func NewUserHashedPassword(hashedPassword []byte) UserHashedPassword {
	return UserHashedPassword(hashedPassword)
}

var (
	ErrInvalidPassword = errors.New("invalid password")
)

func (uhp UserHashedPassword) Compare(password []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(uhp), password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidPassword
	}
	if err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}

	return nil
}
