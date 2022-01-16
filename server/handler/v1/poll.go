package v1

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	openapi "github.com/cs-sysimpl/suzukake/handler/v1/openapi"
	"github.com/cs-sysimpl/suzukake/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Poll struct {
	*Session
	pollService service.Poll
}

func NewPoll(session *Session, pollService service.Poll) *Poll {
	return &Poll{
		Session:     session,
		pollService: pollService,
	}
}

func (p *Poll) PostPolls(c echo.Context) error {
	var poll openapi.PostPollsJSONRequestBody
	err := c.Bind(&poll)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json body")
	}

	session, err := p.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := p.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	title := values.NewPollTitle(string(poll.Title))
	var pollType values.PollType
	switch poll.Type {
	case openapi.PollTypeRadio:
		pollType = values.PollTypeRadio
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid poll type")
	}

	choices := make([]values.ChoiceLabel, 0, len(poll.Question))
	for _, question := range poll.Question {
		choices = append(choices, values.NewChoiceLabel(string(question)))
	}

	var tags []values.TagName
	if poll.Tags == nil {
		tags = []values.TagName{}
	} else {
		tags = make([]values.TagName, 0, len(*poll.Tags))
		for _, tag := range *poll.Tags {
			tags = append(tags, values.NewTagName(string(tag)))
		}
	}

	pollInfo, err := p.pollService.CreatePoll(
		c.Request().Context(),
		user,
		title,
		pollType,
		poll.Deadline,
		choices,
		tags,
	)
	if errors.Is(err, service.ErrNoTag) {
		return echo.NewHTTPError(http.StatusBadRequest, "no tag")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create poll")
	}

	apiPollInfo, err := p.pollInfoToPollSummary(user, pollInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to convert poll info")
	}

	return c.JSON(http.StatusCreated, apiPollInfo)
}

func (p *Poll) GetPolls(c echo.Context, params openapi.GetPollsParams) error {
	session, err := p.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := p.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		user = nil
	}
	if err != nil && !errors.Is(err, ErrNoValue) {
		log.Printf("failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	var searchParams *service.PollSearchParams
	if params.Limit != nil || params.Offset != nil || params.Match != nil {
		searchParams = &service.PollSearchParams{
			Limit:  0,
			Offset: 0,
			Match:  "",
		}

		if params.Limit != nil {
			searchParams.Limit = *params.Limit
		}

		if params.Offset != nil {
			searchParams.Offset = *params.Offset
		}

		if params.Match != nil {
			searchParams.Match = *params.Match
		}
	}

	polls, err := p.pollService.GetPolls(
		c.Request().Context(),
		user,
		searchParams,
	)
	if err != nil {
		log.Printf("failed to get polls: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get polls")
	}

	apiPolls := make([]*openapi.PollSummary, 0, len(polls))
	for _, poll := range polls {
		apiPoll, err := p.pollInfoToPollSummary(user, poll)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to convert poll info")
		}

		apiPolls = append(apiPolls, &apiPoll)
	}

	return c.JSON(http.StatusOK, apiPolls)
}

func (p *Poll) DeletePollsPollID(ctx echo.Context, pollID string) error {
	session, err := p.Session.getSession(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := p.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	uuidPollID, err := uuid.Parse(pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid poll id")
	}

	err = p.pollService.DeletePoll(
		ctx.Request().Context(),
		user,
		values.NewPollIDFromUUID(uuidPollID),
	)
	if errors.Is(err, service.ErrNoPoll) {
		return echo.NewHTTPError(http.StatusBadRequest, "no poll")
	}
	if err != nil {
		log.Printf("failed to get polls: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete poll")
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (p *Poll) GetPollsPollID(c echo.Context, pollID string) error {
	session, err := p.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := p.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		user = nil
	}
	if err != nil && !errors.Is(err, ErrNoValue) {
		log.Printf("failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	uuidPollID, err := uuid.Parse(pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid poll id")
	}

	pollInfo, err := p.pollService.GetPoll(
		c.Request().Context(),
		user,
		values.NewPollIDFromUUID(uuidPollID),
	)
	if errors.Is(err, service.ErrNoPoll) {
		return echo.NewHTTPError(http.StatusBadRequest, "no poll")
	}
	if err != nil {
		log.Printf("failed to get poll: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get poll")
	}

	apiPollInfo, err := p.pollInfoToPollSummary(user, pollInfo)
	if err != nil {
		log.Printf("failed to parse poll info: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to convert poll info")
	}

	return c.JSON(http.StatusOK, apiPollInfo)
}

func (p *Poll) PostPollsClose(c echo.Context, pollID string) error {
	session, err := p.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := p.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	uuidPollID, err := uuid.Parse(pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid poll id")
	}

	err = p.pollService.ClosePoll(
		c.Request().Context(),
		user,
		values.NewPollIDFromUUID(uuidPollID),
	)
	if errors.Is(err, service.ErrNoPoll) {
		return echo.NewHTTPError(http.StatusBadRequest, "no poll")
	}
	if errors.Is(err, service.ErrPollClosed) {
		return echo.NewHTTPError(http.StatusBadRequest, "poll already closed")
	}
	if errors.Is(err, service.ErrNotOwner) {
		return echo.NewHTTPError(http.StatusForbidden, "not owner")
	}
	if err != nil {
		log.Printf("failed to close poll: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to close poll")
	}

	return c.NoContent(http.StatusNoContent)
}

func (p *Poll) pollInfoToPollSummary(user *domain.User, pollInfo *service.PollInfo) (openapi.PollSummary, error) {
	var apiPollType openapi.PollType
	switch pollInfo.Poll.GetPollType() {
	case values.PollTypeRadio:
		apiPollType = openapi.PollTypeRadio
	default:
		return openapi.PollSummary{}, errors.New("invalid poll type")
	}

	var pointerAPITags *[]openapi.PollTag
	if len(pollInfo.Tags) != 0 {
		apiTags := make([]openapi.PollTag, 0, len(pollInfo.Tags))
		for _, tag := range pollInfo.Tags {
			apiTags = append(apiTags, openapi.PollTag{
				Id:   uuid.UUID(tag.GetID()).String(),
				Name: string(tag.GetName()),
			})
		}

		pointerAPITags = &apiTags
	}

	apiChoices := make(openapi.Questions, 0, len(pollInfo.Choices))
	for _, choice := range pollInfo.Choices {
		apiChoices = append(apiChoices, openapi.Choice{
			Id:     uuid.UUID(choice.GetID()).String(),
			Choice: string(choice.GetLabel()),
		})
	}

	deadline := pollInfo.GetDeadline()
	var apiDeadline *time.Time
	if deadline.Valid {
		apiDeadline = &deadline.Time
	}

	var apiPollStatus openapi.PollStatus
	switch {
	case pollInfo.Poll.IsExpired():
		apiPollStatus = openapi.PollStatusOutdated
	case pollInfo.Poll.GetDeadline().Valid:
		apiPollStatus = openapi.PollStatusLimited
	default:
		apiPollStatus = openapi.PollStatusOpened
	}

	var apiUserStatus openapi.UserStatus
	switch {
	case pollInfo.Owner.GetID() == user.GetID():
		apiUserStatus = openapi.UserStatus{
			AccsessMode: openapi.UserStatusAccsessModeCanAsccessDetails,
			IsOwner:     true,
		}
	case pollInfo.Poll.IsExpired():
		apiUserStatus = openapi.UserStatus{
			AccsessMode: openapi.UserStatusAccsessModeCanAsccessDetails,
			IsOwner:     false,
		}
	case user == nil:
		apiUserStatus = openapi.UserStatus{
			AccsessMode: openapi.UserStatusAccsessModeOnlyBrowsable,
			IsOwner:     false,
		}
	case pollInfo.Response == nil:
		apiUserStatus = openapi.UserStatus{
			AccsessMode: openapi.UserStatusAccsessModeCanAnswer,
			IsOwner:     false,
		}
	default:
		apiUserStatus = openapi.UserStatus{
			AccsessMode: openapi.UserStatusAccsessModeCanAsccessDetails,
			IsOwner:     false,
		}
	}

	apiPollInfo := openapi.PollSummary{
		PollId: openapi.PollID(uuid.UUID(pollInfo.Poll.GetID()).String()),
		PollBase: openapi.PollBase{
			Title:    string(pollInfo.Poll.GetTitle()),
			Type:     apiPollType,
			Tags:     pointerAPITags,
			Question: apiChoices,
			Deadline: apiDeadline,
		},
		CreatedAt: pollInfo.GetCreatedAt(),
		Owner: openapi.User{
			Uuid: uuid.UUID(pollInfo.Owner.GetID()).String(),
			Name: openapi.UserName(pollInfo.Owner.GetName()),
		},
		QStatus:    apiPollStatus,
		UserStatus: apiUserStatus,
	}

	return apiPollInfo, nil
}
