package v1

import (
	"errors"
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
