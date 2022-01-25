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

// CanRead userはnullable
func (p *PollAuthority) CanRead(ctx context.Context, user *domain.User, owner *domain.User, poll *domain.Poll) (bool, error) {
	if poll.IsExpired() || owner.GetID() == user.GetID() {
		return true, nil
	}

	if user == nil {
		return false, nil
	}

	_, err := p.responseRepository.GetResponseByUserIDAndPollID(ctx, user.GetID(), poll.GetID(), repository.LockTypeNone)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get response: %w", err)
	}

	return true, nil
}

func (p *PollAuthority) CanResponse(ctx context.Context, user *domain.User, owner *domain.User, poll *domain.Poll) (bool, error) {
	if poll.IsExpired() || owner.GetID() == user.GetID() {
		return false, nil
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
