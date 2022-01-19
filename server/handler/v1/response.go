package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/cs-sysimpl/suzukake/domain/values"
	openapi "github.com/cs-sysimpl/suzukake/handler/v1/openapi"
	"github.com/cs-sysimpl/suzukake/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Response struct {
	*Session
	responseService service.Response
}

func NewResponse(session *Session, responseService service.Response) *Response {
	return &Response{
		Session:         session,
		responseService: responseService,
	}
}

func (r *Response) PostPollsPollID(c echo.Context, pollID string) error {
	session, err := r.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := r.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		user = nil
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	var req openapi.PostPollsPollIDJSONRequestBody
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	choiceIDs := make([]values.ChoiceID, 0, len(req.Answer))
	for _, answer := range req.Answer {
		uuidChoiceID, err := uuid.Parse(answer)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid choice id")
		}

		choiceIDs = append(choiceIDs, values.NewChoiceIDFromUUID(uuidChoiceID))
	}

	uuidPollID, err := uuid.Parse(pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid poll id")
	}

	response, err := r.responseService.CreateResponse(
		c.Request().Context(),
		user,
		values.NewPollIDFromUUID(uuidPollID),
		choiceIDs,
		values.NewCommentContent(req.Comment),
	)
	if errors.Is(err, service.ErrNoPoll) {
		return echo.NewHTTPError(http.StatusBadRequest, "poll not found")
	}
	if errors.Is(err, service.ErrResponseAlreadyExists) {
		return echo.NewHTTPError(http.StatusBadRequest, "response already exists")
	}
	if errors.Is(err, service.ErrTooManyChoice) {
		return echo.NewHTTPError(http.StatusBadRequest, "too many choice")
	}
	if errors.Is(err, service.ErrNoChoice) {
		return echo.NewHTTPError(http.StatusBadRequest, "choice not found")
	}
	if errors.Is(err, service.ErrDuplicateChoices) {
		return echo.NewHTTPError(http.StatusBadRequest, "duplicate choices")
	}
	if err != nil {
		log.Printf("failed to post response: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to post response")
	}

	choices := make([]string, 0, len(response.Choices))
	for _, choice := range response.Choices {
		choices = append(choices, uuid.UUID(choice.GetID()).String())
	}

	var comment *string
	if response.Comment != nil {
		commentContent := string(response.Comment.GetContent())
		comment = &commentContent
	}

	return c.JSON(http.StatusCreated, openapi.Response{
		Answer:    choices,
		Comment:   comment,
		CreatedAt: response.Response.GetCreatedAt(),
	})
}
