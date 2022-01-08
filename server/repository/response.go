package repository

import (
	"context"
	"time"

	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Response interface {
	GetIdsByPollId(ctx context.Context, pollID values.PollID) ([]*values.ResponseID, error)
	GetCreatedAtById(ctx context.Context, id values.ResponseID) (time.Time, error)
	GetUserIdAtById(ctx context.Context, id values.ResponseID) (values.UserID, error)
}
