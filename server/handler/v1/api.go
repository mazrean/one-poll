package v1

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config ./openapi/config.yaml ../../../docs/openapi/openapi.yaml
//go:generate go fmt ./openapi/openapi.gen.go

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	openapi "github.com/mazrean/one-poll/handler/v1/openapi"
	oapiMiddleware "github.com/oapi-codegen/echo-middleware"
)

type unimplemented interface {
	// (POST /tags)
	PostTags(ctx echo.Context) error
	// webauthnの登録情報削除
	// (DELETE /webauthn/credentials)
	DeleteWebauthnCredentials(ctx echo.Context) error
}

type API struct {
	*Checker
	*User
	*Poll
	*Tag
	*Comment
	*Response
	*WebAuthn
	unimplemented
}

func NewAPI(
	checker *Checker,
	user *User,
	poll *Poll,
	tag *Tag,
	comment *Comment,
	response *Response,
	webAuthn *WebAuthn,
) *API {
	return &API{
		Checker:  checker,
		User:     user,
		Poll:     poll,
		Tag:      tag,
		Comment:  comment,
		Response: response,
		WebAuthn: webAuthn,
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
