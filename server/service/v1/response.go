package v1

import (
	"github.com/cs-sysimpl/suzukake/repository"
)

type Response struct {
	db                 repository.DB
	pollRepository     repository.Poll
	responseRepository repository.Response
	commentRepository  repository.Comment
	choiceRepository   repository.Choice
	pollAuthority      *PollAuthority
}

func NewResponse(
	db repository.DB,
	pollRepository repository.Poll,
	responseRepository repository.Response,
	commentRepository repository.Comment,
	choiceRepository repository.Choice,
	pollAuthority *PollAuthority,
) *Response {
	return &Response{
		db:                 db,
		pollRepository:     pollRepository,
		responseRepository: responseRepository,
		commentRepository:  commentRepository,
		choiceRepository:   choiceRepository,
		pollAuthority:      pollAuthority,
	}
}
