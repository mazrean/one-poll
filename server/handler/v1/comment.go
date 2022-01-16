package v1

import (
	"log"
	"net/http"

	"github.com/cs-sysimpl/suzukake/domain/values"
	openapi "github.com/cs-sysimpl/suzukake/handler/v1/openapi"
	"github.com/cs-sysimpl/suzukake/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Comment struct {
	commentService service.Comment
}

func NewComment(
	commentService service.Comment,
) *Comment {
	return &Comment{
		commentService: commentService,
	}
}

// todo Comment の数の制限については未実装
func (c Comment) GetPollsPollIDComments(ctx echo.Context, pollID string, params openapi.GetPollsPollIDCommentsParams) error {
	uuidPollID, err := uuid.Parse(pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid poll id")
	}
	commentInfos, err := c.commentService.GetComments(
		ctx.Request().Context(),
		values.NewPollIDFromUUID(uuidPollID),
	)
	if err != nil {
		log.Printf("failed to get comments: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get comments")
	}

	comments := make([]openapi.PollComment, 0, len(commentInfos))

	for _, commentInfo := range commentInfos {
		comment := openapi.PollComment{
			Content:   string(commentInfo.GetContent()),
			CreatedAt: commentInfo.GetCreatedAt(),
			User: openapi.User{
				Name: openapi.UserName(commentInfo.CommentUser.GetName()),
				Uuid: uuid.UUID(commentInfo.CommentUser.GetID()).String(),
			},
		}
		comments = append(comments, comment)
	}

	return ctx.JSON(http.StatusOK, comments)
}
