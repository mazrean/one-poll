package repository

import (
	"context"
	"time"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Response interface {
	GetResponsesByPollID(ctx context.Context, pollID values.PollID) ([]*domain.Response, error)
	GetResponseByUserAndPollID(ctx context.Context, userID values.UserID, pollID values.PollID, lockType LockType) (*domain.Response, error)
	GetIDsByPollID(ctx context.Context, pollID values.PollID) ([]*values.ResponseID, error)
	GetCreatedAtByID(ctx context.Context, id values.ResponseID) (time.Time, error)
	GetUserIDByID(ctx context.Context, id values.ResponseID) (values.UserID, error)
}

type ResponseInfo struct {
	*domain.Response
	Respondent *domain.User
}
