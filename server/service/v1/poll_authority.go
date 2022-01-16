package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
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

func (p *PollAuthority) CanRead(ctx context.Context, user *domain.User, poll *domain.Poll) (bool, error) {
	if poll.IsExpired() {
		return true, nil
	}

	_, err := p.responseRepository.GetResponseByUserIDAndPollID(ctx, user.GetID(), poll.GetID(), repository.LockTypeNone)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get response: %w", err)
	}

	return false, nil
}
