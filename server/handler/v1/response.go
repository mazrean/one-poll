package v1

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/one-poll/domain/values"
	openapi "github.com/mazrean/one-poll/handler/v1/openapi"
	"github.com/mazrean/one-poll/service"
	"github.com/oapi-codegen/runtime/types"
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
	for _, uuidChoiceID := range req.Answer {
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

	choices := make([]types.UUID, 0, len(response.Choices))
	for _, choice := range response.Choices {
		choices = append(choices, types.UUID(choice.GetID()))
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

func (r *Response) GetPollsPollIDResults(c echo.Context, pollID string) error {
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

	uuidPollID, err := uuid.Parse(pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid poll id")
	}

	result, err := r.responseService.GetResult(
		c.Request().Context(),
		user,
		values.NewPollIDFromUUID(uuidPollID),
	)
	if errors.Is(err, service.ErrNoPoll) {
		return echo.NewHTTPError(http.StatusBadRequest, "poll not found")
	}
	if errors.Is(err, service.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	if err != nil {
		log.Printf("failed to get poll: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get poll")
	}

	var polltype openapi.PollType
	switch result.Poll.GetPollType() {
	case values.PollTypeRadio:
		polltype = openapi.Radio
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid poll type")
	}

	results := make([]openapi.Result, 0, len(result.Items))
	for _, item := range result.Items {
		results = append(results, openapi.Result{
			Id:     types.UUID(item.Choice.GetID()),
			Choice: string(item.GetLabel()),
			Count:  item.Count,
		})
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	if result.GetDeadline().Valid && result.GetDeadline().Time.Before(time.Now()) {
		c.Response().Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=31536000, stale-if-error=31536000")
	}

	return c.JSON(http.StatusOK, openapi.PollResults{
		PollId: openapi.PollID(result.GetID()),
		Type:   polltype,
		Count:  result.Count,
		Result: results,
	})
}
