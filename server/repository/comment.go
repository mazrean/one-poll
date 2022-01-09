package repository

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Comment interface {
	CreateComment(ctx context.Context, comment *domain.Comment, pollID values.PollID) error
	GetCommentByPollID(ctx context.Context, pollID values.PollID) (*domain.Comment, error)
}
