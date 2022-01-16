package v1

import "github.com/cs-sysimpl/suzukake/service"

type Poll struct {
	*Session
	pollService service.Poll
}

func NewPoll(session *Session, pollService service.Poll) *Poll {
	return &Poll{
		Session:     session,
		pollService: pollService,
	}
}
