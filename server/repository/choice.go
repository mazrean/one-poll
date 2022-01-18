package repository

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Choice interface {
	CreateChoices(ctx context.Context, pollID values.PollID, choices []*domain.Choice) error
	GetChoices(ctx context.Context, ids []values.ChoiceID, lockType LockType) ([]*domain.Choice, error)
	GetChoicesByPollIDs(ctx context.Context, pollIDs []values.PollID, lockType LockType) (map[values.PollID][]*domain.Choice, error)
	GetChoicesByPollID(ctx context.Context, pollID values.PollID, lockType LockType) ([]*domain.Choice, error)
}
