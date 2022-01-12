package repository

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Comment interface {
	CreateComment(ctx context.Context, pollID values.PollID, comment *domain.Comment) error
	GetCommentsByResponseIDs(ctx context.Context, responseIDs values.ResponseID) (map[values.ResponseID]*domain.Comment, error)
}
