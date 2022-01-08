package service

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Comment interface {
	GetComments(ctx context.Context, pollID values.PollID) (CommentsService, error)
}

type CommentService struct {
	domain.Response
	domain.Comment
	commentUser domain.User
}

type CommentsService struct {
	values.PollID
	comments []CommentService
}
