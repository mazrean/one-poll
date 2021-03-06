package v1

import (
	"log"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/labstack/echo/v4"
	openapi "github.com/mazrean/one-poll/handler/v1/openapi"
	"github.com/mazrean/one-poll/service"
)

type Tag struct {
	tagService service.Tag
}

func NewTag(service service.Tag) *Tag {
	return &Tag{
		tagService: service,
	}
}

func (t *Tag) GetTags(c echo.Context) error {
	tags, err := t.tagService.GetTags(c.Request().Context())
	if err != nil {
		log.Printf("failed to get tags: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get tags")
	}

	apiTags := make([]openapi.PollTag, 0, len(tags))
	for _, tag := range tags {
		apiTags = append(apiTags, openapi.PollTag{
			Id:   types.UUID(tag.GetID()),
			Name: string(tag.GetName()),
		})
	}

	return c.JSON(http.StatusOK, apiTags)
}
