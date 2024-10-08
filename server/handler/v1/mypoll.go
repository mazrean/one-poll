package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	openapi "github.com/mazrean/one-poll/handler/v1/openapi"
)

func (p *Poll) GetUsersMeOwners(ctx echo.Context) error {
	session, err := p.Session.getSession(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := p.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil && !errors.Is(err, ErrNoValue) {
		log.Printf("failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	pollInfos, err := p.pollService.GetOwnerPolls(ctx.Request().Context(), user)
	if err != nil {
		log.Printf("failed to get owner polls: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get owner polls")
	}

	apiPolls := make([]*openapi.PollSummary, 0, len(pollInfos))

	for _, pollInfo := range pollInfos {
		apiPoll, err := p.pollInfoToPollSummary(user, pollInfo)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to convert poll info")
		}

		apiPolls = append(apiPolls, &apiPoll)
	}

	ctx.Response().Header().Set("Cache-Control", "no-store")
	return ctx.JSON(http.StatusOK, apiPolls)
}

func (p *Poll) GetUsersMeAnswers(ctx echo.Context) error {
	session, err := p.Session.getSession(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := p.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil && !errors.Is(err, ErrNoValue) {
		log.Printf("failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	pollInfos, err := p.pollService.GetAnsweredPolls(ctx.Request().Context(), user)
	if err != nil {
		log.Printf("failed to get answered polls: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get answered polls")
	}

	apiPolls := make([]*openapi.PollSummary, 0, len(pollInfos))

	for _, pollInfo := range pollInfos {
		apiPoll, err := p.pollInfoToPollSummary(user, pollInfo)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to convert poll info")
		}

		apiPolls = append(apiPolls, &apiPoll)
	}

	ctx.Response().Header().Set("Cache-Control", "no-store")
	return ctx.JSON(http.StatusOK, apiPolls)
}
