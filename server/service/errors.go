package service

import "errors"

var (
	ErrNoPoll                = errors.New("no poll")
	ErrNoTag                 = errors.New("no tag")
	ErrNoChoice              = errors.New("no choice")
	ErrTooManyChoice         = errors.New("too many choice")
	ErrDuplicateChoices      = errors.New("duplicate choices")
	ErrNotOwner              = errors.New("not owner")
	ErrPollClosed            = errors.New("poll closed")
	ErrResponseAlreadyExists = errors.New("response already exists")
	ErrForbidden             = errors.New("forbidden")
)
