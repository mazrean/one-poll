package repository

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Comment interface {
	CreateComment(ctx context.Context, responseID values.ResponseID, comment *domain.Comment) error
	GetCommentsByResponseIDs(ctx context.Context, responseIDs []values.ResponseID, params CommentGetParams) (map[values.ResponseID]*domain.Comment, error)
}

type CommentGetParams struct {
	Limit  *int
	Offset *int
}
