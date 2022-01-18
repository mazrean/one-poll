package service

import "errors"

var (
	ErrNoPoll           = errors.New("no poll")
	ErrNoTag            = errors.New("no tag")
	ErrDuplicateChoices = errors.New("duplicate choices")
	ErrNotOwner         = errors.New("not owner")
	ErrPollClosed       = errors.New("poll closed")
)
