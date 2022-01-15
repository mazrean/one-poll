package repository

import (
	"context"
	"database/sql"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Poll interface {
	CreatePoll(ctx context.Context, poll *domain.Poll, owner values.UserID) error
	UpdatePollDeadline(ctx context.Context, id values.PollID, deadline sql.NullTime) error
	GetPolls(ctx context.Context, params *PollSearchParams) ([]*PollInfo, error)
	GetPoll(ctx context.Context, id values.PollID, lockType LockType) (*PollInfo, error)
	DeletePoll(ctx context.Context, id values.PollID) error
	AddTags(ctx context.Context, pollID values.PollID, tags []values.TagID) error
}

type PollInfo struct {
	*domain.Poll
	Owner *domain.User
}

type PollSearchParams struct {
	Limit  int
	Offset int
	Match  string
}
