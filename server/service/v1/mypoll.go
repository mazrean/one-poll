package v1

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/service"
)

func (p *Poll) GetOwnerPolls(ctx context.Context, owner *domain.User) ([]*service.PollInfo, error) {
	// TODO Implementation
	pollInfos := make([]*service.PollInfo, 0, 0)
	return pollInfos, nil
}
func (p *Poll) GetAnsweredPolls(ctx context.Context, owner *domain.User) ([]*service.PollInfo, error) {
	// TODO Implementation
	pollInfos := make([]*service.PollInfo, 0, 0)
	return pollInfos, nil
}
