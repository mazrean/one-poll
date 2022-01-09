package values

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type (
	PollID    uuid.UUID
	PollTitle string
	PollType  int8
)

func NewPollID() PollID {
	return PollID(uuid.New())
}

func NewPollIDFromUUID(id uuid.UUID) PollID {
	return PollID(id)
}

func NewPollTitle(title string) PollTitle {
	return PollTitle(title)
}

var (
	ErrPollTitleEmpty   = errors.New("poll title is empty")
	ErrPollTitleTooLong = errors.New("poll title is too long")
)

func (p PollTitle) Validate() error {
	if len(p) == 0 {
		return ErrPollTitleEmpty
	}

	if utf8.RuneCountInString(string(p)) > 50 {
		return ErrPollTitleTooLong
	}

	return nil
}

const (
	PollTypeRadio = iota + 1
)
