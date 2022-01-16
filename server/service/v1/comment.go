package v1

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/repository"
	"github.com/cs-sysimpl/suzukake/service"
)

type Comment struct {
	responseRepository repository.Response
	commentRepository  repository.Comment
}

func (c *Comment) GetComments(ctx context.Context, pollID values.PollID) ([]service.CommentInfo, error) {
	responseInfos, err := c.responseRepository.GetResponsesByPollID(ctx, pollID)
	if err != nil {
		return nil, fmt.Errorf("failed to get responseInfos: %w", err)
	}
	responseIDs := make([]values.ResponseID, 0, len(responseInfos))
	for _, responseInfo := range responseInfos {
		responseIDs = append(responseIDs, responseInfo.Response.GetID())
	}
	dbComments, err := c.commentRepository.GetCommentsByResponseIDs(ctx, responseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	commentInfos := make([]service.CommentInfo, 0, len(responseInfos))
	for _, responseInfo := range responseInfos {
		commentInfo := service.CommentInfo{
			Response:    *responseInfo.Response,
			Comment:     *dbComments[responseInfo.Response.GetID()],
			CommentUser: *responseInfo.Respondent}
		commentInfos = append(commentInfos, commentInfo)
	}

	return commentInfos, nil

}
