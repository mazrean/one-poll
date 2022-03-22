package repository

import (
	"context"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
)

type Response interface {
	CreateResponse(ctx context.Context, userID values.UserID, pollID values.PollID, response *domain.Response, choiceIDs []values.ChoiceID) error
	GetResponsesByPollID(ctx context.Context, pollID values.PollID) ([]*ResponseInfo, error)
	GetResponseByUserIDAndPollID(ctx context.Context, userID values.UserID, pollID values.PollID, lockType LockType) (*domain.Response, error)
	GetResponsesByUserIDAndPollIDs(ctx context.Context, userID values.UserID, pollIDs []values.PollID, lockType LockType) (map[values.PollID]*domain.Response, error)
	GetStatistics(ctx context.Context, pollID values.PollID) (*Statistics, error)
}

type ResponseInfo struct {
	*domain.Response
	Respondent *domain.User
}

type Statistics struct {
	Count       int
	ChoiceCount map[values.ChoiceID]int
}
