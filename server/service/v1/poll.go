package v1

import (
	"github.com/cs-sysimpl/suzukake/repository"
)

type Poll struct {
	db               repository.DB
	pollRepository   repository.Poll
	choiceRepository repository.Choice
	tagRepository    repository.Tag
}

func NewPoll(
	db repository.DB,
	pollRepository repository.Poll,
	choiceRepository repository.Choice,
	tagRepository repository.Tag,
) *Poll {
	return &Poll{
		db:               db,
		pollRepository:   pollRepository,
		choiceRepository: choiceRepository,
		tagRepository:    tagRepository,
	}
}
