package v1

import (
	"net/http"
	"time"

	openapi "github.com/cs-sysimpl/suzukake/handler/v1/openapi"
	"github.com/labstack/echo/v4"
)

type Comment struct {
}

// [TODO] service の GetComments を使う
func (c Comment) GetPollsPollIDComments(ctx echo.Context, pollID string, params openapi.GetPollsPollIDCommentsParams) error {

	return ctx.JSON(http.StatusOK, []openapi.PollComment{
		{
			Content:   "content",
			CreatedAt: time.Now(),
			User:      openapi.User{Name: "hoge", Uuid: "uuid"},
		},
	})
}
