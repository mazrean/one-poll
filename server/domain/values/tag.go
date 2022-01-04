package values

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type (
	TagID   uuid.UUID
	TagName string
)

func NewTagID() TagID {
	return TagID(uuid.New())
}

func NewTagIDFromUUID(id uuid.UUID) TagID {
	return TagID(id)
}

func NewTagName(name string) TagName {
	return TagName(name)
}

var (
	ErrTagNameEmpty   = errors.New("tag name is empty")
	ErrTagNameTooLong = errors.New("tag name is too long")
)

func (tn TagName) Validate() error {
	if len(tn) == 0 {
		return ErrTagNameEmpty
	}

	if utf8.RuneCountInString(string(tn)) > 50 {
		return ErrTagNameTooLong
	}

	return nil
}
