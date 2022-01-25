package v1

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/repository"
	"github.com/cs-sysimpl/suzukake/service"
)

type Comment struct {
	responseRepository repository.Response
	commentRepository  repository.Comment
	pollRepository     repository.Poll
	pollAuthority      *PollAuthority
}

func NewComment(
	responseRepository repository.Response,
	commentRepository repository.Comment,
	pollRepository repository.Poll,
	pollAuthority *PollAuthority,
) *Comment {
	return &Comment{
		responseRepository: responseRepository,
		commentRepository:  commentRepository,
		pollRepository:     pollRepository,
		pollAuthority:      pollAuthority,
	}
}

func (c *Comment) GetComments(ctx context.Context, pollID values.PollID, user *domain.User, params service.CommentGetParams) ([]service.CommentInfo, error) {
	pollInfo, err := c.pollRepository.GetPoll(ctx, pollID, repository.LockTypeNone)
	if err != nil {
		return nil, fmt.Errorf("failed to get polls: %w", err)
	}

	canRead, err := c.pollAuthority.CanRead(ctx, user, pollInfo.Owner, pollInfo.Poll)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	if !canRead {
		return nil, service.ErrForbidden
	}

	responseInfos, err := c.responseRepository.GetResponsesByPollID(ctx, pollID)
	if err != nil {
		return nil, fmt.Errorf("failed to get responseInfos: %w", err)
	}
	responseIDs := make([]values.ResponseID, 0, len(responseInfos))
	for _, responseInfo := range responseInfos {
		responseIDs = append(responseIDs, responseInfo.Response.GetID())
	}
	// コメントの存在するResponseIDに対してのみキーとしたマップを作成する。
	dbComments, err := c.commentRepository.GetCommentsByResponseIDs(ctx, responseIDs, (repository.CommentGetParams)(params))
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	commentInfos := make([]service.CommentInfo, 0, len(responseInfos))
	for _, responseInfo := range responseInfos {
		comment, exists := dbComments[responseInfo.Response.GetID()]
		commentInfo := service.CommentInfo{
			Response:    *responseInfo.Response,
			CommentUser: *responseInfo.Respondent}
		if exists {
			commentInfo.Comment = *comment
		}

		commentInfos = append(commentInfos, commentInfo)
	}

	return commentInfos, nil

}
