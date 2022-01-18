package service

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Comment interface {
	GetComments(ctx context.Context, pollID values.PollID, user *domain.User) ([]CommentInfo, error)
}

type CommentInfo struct {
	domain.Response
	domain.Comment
	CommentUser domain.User
}
