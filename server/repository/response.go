package repository

import (
	"context"
	"time"

	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Response interface {
	GetIDsByPollID(ctx context.Context, pollID values.PollID) ([]*values.ResponseID, error)
	GetCreatedAtByID(ctx context.Context, id values.ResponseID) (time.Time, error)
	GetUserIDByID(ctx context.Context, id values.ResponseID) (values.UserID, error)
}
