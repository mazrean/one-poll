package values

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type (
	ChoiceID    uuid.UUID
	ChoiceLabel string
)

func NewChoiceID() ChoiceID {
	return ChoiceID(uuid.New())
}

func NewChoiceIDFromUUID(id uuid.UUID) ChoiceID {
	return ChoiceID(id)
}

func NewChoiceLabel(label string) ChoiceLabel {
	return ChoiceLabel(label)
}

var (
	ErrChoiceLabelEmpty   = errors.New("choice label is empty")
	ErrChoiceLabelTooLong = errors.New("choice label is too long")
)

func (cl ChoiceLabel) Validate() error {
	if len(cl) == 0 {
		return ErrPollTitleEmpty
	}

	if utf8.RuneCountInString(string(cl)) > 50 {
		return ErrPollTitleTooLong
	}

	return nil
}
