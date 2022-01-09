package values

import (
	"errors"
	"fmt"
	"unicode/utf8"

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

var (
	ErrUserNameTooShort    = errors.New("user name is too short")
	ErrUserNameTooLong     = errors.New("user name is too long")
	ErrUserNameInvalidRune = errors.New("user name contains invalid rune")
)

func (un UserName) Validate() error {
	length := utf8.RuneCountInString(string(un))
	if length < 4 {
		return ErrUserNameTooShort
	}

	if length > 16 {
		return ErrUserNameTooLong
	}

	for _, r := range un {
		if !('0' <= r && r <= '9') && !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') && r != '_' {
			return ErrUserNameInvalidRune
		}
	}

	return nil
}

func NewUserPassword(password []byte) UserPassword {
	return UserPassword(password)
}

var (
	ErrUserPasswordTooShort    = errors.New("user password is too short")
	ErrUserPasswordTooLong     = errors.New("user password is too long")
	ErrUserPasswordInvalidRune = errors.New("user password contains invalid rune")
)

func (up UserPassword) Validate() error {
	length := utf8.RuneCountInString(string(up))
	if length < 8 {
		return ErrUserPasswordTooShort
	}
	if length > 50 {
		return ErrUserPasswordTooLong
	}

	for _, r := range up {
		if !('0' <= r && r <= '9') && !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') {
			return ErrUserPasswordInvalidRune
		}
	}

	return nil
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
