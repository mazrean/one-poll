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
	// GetResult userがnullableであることに注意
	GetResult(ctx context.Context, user *domain.User, pollID values.PollID) (*Result, error)
}

type ResponseInfo struct {
	*domain.Poll
	*domain.Response
	// Comment nullableなことに注意
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
