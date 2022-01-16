package service

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Response interface {
	CreateResponse(
		ctx context.Context,
		user *domain.User,
		pollID values.PollID,
		choices []values.ChoiceID,
		comment values.CommentContent,
	) (*ResponseInfo, error)
	GetResult(ctx context.Context, user *domain.User, pollID values.PollID) (*Result, error)
}

type ResponseInfo struct {
	*domain.Poll
	*domain.Response
	*domain.Comment
	Choices []*domain.Choice
}

type Result struct {
	*domain.Poll
	Count int
	Items []*ResultItem
}

type ResultItem struct {
	*domain.Choice
	Count int
}
