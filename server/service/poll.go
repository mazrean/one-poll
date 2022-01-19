package service

import (
	"context"
	"time"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Poll interface {
	CreatePoll(
		ctx context.Context,
		user *domain.User,
		title values.PollTitle,
		pollType values.PollType,
		deadline *time.Time,
		choices []values.ChoiceLabel,
		tags []values.TagName,
	) (*PollInfo, error)
	GetPolls(ctx context.Context, user *domain.User, params *PollSearchParams) ([]*PollInfo, error)
	GetOwnerPolls(ctx context.Context, owner *domain.User) ([]*PollInfo, error)
	GetAnsweredPolls(ctx context.Context, owner *domain.User) ([]*PollInfo, error)
	GetPoll(ctx context.Context, user *domain.User, id values.PollID) (*PollInfo, error)
	ClosePoll(ctx context.Context, user *domain.User, id values.PollID) error
	DeletePoll(ctx context.Context, user *domain.User, id values.PollID) error
}

type PollInfo struct {
	*domain.Poll
	Tags    []*domain.Tag
	Choices []*domain.Choice
	Owner   *domain.User
	// Response nullableなことに注意
	Response *domain.Response
}

type PollSearchParams struct {
	Limit  int
	Offset int
	Match  string
}
