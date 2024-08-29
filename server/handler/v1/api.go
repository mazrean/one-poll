package v1

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config ./openapi/config.yaml ../../../docs/openapi/openapi.yaml
//go:generate go fmt ./openapi/openapi.gen.go

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	openapi "github.com/mazrean/one-poll/handler/v1/openapi"
	oapiMiddleware "github.com/oapi-codegen/echo-middleware"
)

//go:embed static/*
var staticFS embed.FS

type unimplemented interface {
	// (POST /tags)
	PostTags(ctx echo.Context) error
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
	e.Use(middleware.RewriteWithConfig(middleware.RewriteConfig{
		Rules: map[string]string{
			"/":         "/index.html",
			"/signup":   "/index.html",
			"/signup/":  "/index.html",
			"/singin":   "/index.html",
			"/signin/":  "/index.html",
			"/profile":  "/index.html",
			"/profile/": "/index.html",
		},
		RegexRules: map[*regexp.Regexp]string{
			regexp.MustCompile("^/details/(.*)$"): "/index.html",
		},
	}))
	e.Use(StaticBrotliMiddleware)

	staticFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		return fmt.Errorf("failed to get staticFS: %w", err)
	}
	e.StaticFS("/", staticFS)

	api := e.Group("/api", oapiMiddleware.OapiRequestValidatorWithOptions(swagger, &oapiMiddleware.Options{
		Options: openapi3filter.Options{
			MultiError:         true,
			AuthenticationFunc: a.Checker.check,
		},
	}))
	openapi.RegisterHandlersWithBaseURL(api, a, "")

	return e.Start(addr)
}

func StaticBrotliMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL.Path)
		if c.Request().Method != "GET" ||
			strings.HasPrefix(c.Request().URL.Path, "/api") {
			return next(c)
		}

		exts := []string{".html", ".css", ".js"}
		for _, ext := range exts {
			if strings.HasSuffix(c.Request().URL.Path, ext) {
				goto AFTER_EXT_LOOP
			}
		}
		return next(c)

	AFTER_EXT_LOOP:
		acceptEncodings := strings.Split(strings.TrimSpace(c.Request().Header.Get("Accept-Encoding")), ",")
		for _, acceptEncoding := range acceptEncodings {
			encoding, _, _ := strings.Cut(acceptEncoding, ";")
			if encoding == "br" {
				goto AFTER_ACCEPT_ENCODING_LOOP
			}
		}
		return next(c)

	AFTER_ACCEPT_ENCODING_LOOP:
		c.Response().Header().Set("Content-Encoding", "br")
		c.Request().URL.Path += ".br"

		return next(c)
	}
}
