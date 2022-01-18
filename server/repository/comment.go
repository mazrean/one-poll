package repository

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Comment interface {
	CreateComment(ctx context.Context, responseID values.ResponseID, comment *domain.Comment) error
	GetCommentsByResponseIDs(ctx context.Context, responseIDs []values.ResponseID) (map[values.ResponseID]*domain.Comment, error)
}
