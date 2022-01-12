package v1

import (
	"github.com/cs-sysimpl/suzukake/repository"
)

type PollAuthority struct {
	responseRepository repository.Response
}

func NewPollAuthority(responseRepository repository.Response) *PollAuthority {
	return &PollAuthority{
		responseRepository: responseRepository,
	}
}
