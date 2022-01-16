package repository

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Response interface {
	GetResponsesByPollID(ctx context.Context, pollID values.PollID) ([]*domain.Response, error)
	GetResponseByUserIDAndPollID(ctx context.Context, userID values.UserID, pollID values.PollID, lockType LockType) (*domain.Response, error)
	GetResponsesByUserIDAndPollIDs(ctx context.Context, userID values.UserID, pollIDs []values.PollID, lockType LockType) (map[values.PollID]*domain.Response, error)
}

type ResponseInfo struct {
	*domain.Response
	Respondent *domain.User
}
