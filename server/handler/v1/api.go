package v1

//go:generate sh -c "go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest -generate types,server,spec ../../../docs/openapi/openapi.yaml > openapi/openapi.gen.go"
//go:generate go fmt ./openapi/openapi.gen.go

import (
	"fmt"

	openapi "github.com/cs-sysimpl/suzukake/handler/v1/openapi"
	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type unimplemented interface {
	// (GET /tags)
	GetTags(ctx echo.Context) error

	// (POST /tags)
	PostTags(ctx echo.Context) error

	// (GET /users/me/answers)
	GetUsersMeAnswers(ctx echo.Context) error

	// (GET /users/me/owners)
	GetUsersMeOwners(ctx echo.Context) error
}

type API struct {
	*Checker
	*User
	*Poll
	*Comment
	*Response
	unimplemented
}

func NewAPI(
	checker *Checker,
	user *User,
	poll *Poll,
	comment *Comment,
	response *Response,
) *API {
	return &API{
		Checker:  checker,
		User:     user,
		Poll:     poll,
		Comment:  comment,
		Response: response,
	}
}

func (a *API) Start(addr string) error {
	e := echo.New()

	swagger, err := openapi.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to get openapi: %w", err)
	}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, &oapiMiddleware.Options{
		Options: openapi3filter.Options{
			MultiError:         true,
			AuthenticationFunc: a.Checker.check,
		},
	}))

	openapi.RegisterHandlersWithBaseURL(e, a, "/api")

	return e.Start(addr)
}
