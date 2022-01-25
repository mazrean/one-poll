package repository

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Tag interface {
	CreateTag(ctx context.Context, tag *domain.Tag) error
	GetTags(ctx context.Context) ([]*domain.Tag, error)
	GetTagsByName(ctx context.Context, names []values.TagName, lockType LockType) ([]*domain.Tag, error)
	GetTagsByPollIDs(ctx context.Context, pollIDs []values.PollID, lockType LockType) (map[values.PollID][]*domain.Tag, error)
	GetTagsByPollID(ctx context.Context, pollID values.PollID, lockType LockType) ([]*domain.Tag, error)
}
