package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/cs-sysimpl/suzukake/domain/values"
	openapi "github.com/cs-sysimpl/suzukake/handler/v1/openapi"
	"github.com/cs-sysimpl/suzukake/service"
	"github.com/labstack/echo/v4"
)

type User struct {
	*Session
	authorizationService service.Authorization
}

func NewUser(
	session *Session,
	authorizationService service.Authorization,
) *User {
	return &User{
		Session:              session,
		authorizationService: authorizationService,
	}
}

func (u *User) PostUsers(c echo.Context) error {
	var userRequest openapi.PostUsersJSONRequestBody
	err := c.Bind(&userRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json body")
	}

	name := values.NewUserName(string(userRequest.Name))
	err = name.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user name")
	}

	password := values.NewUserPassword([]byte(userRequest.Password))
	err = password.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user password")
	}

	_, err = u.authorizationService.Signup(
		c.Request().Context(),
		name,
		password,
	)
	if errors.Is(err, service.ErrUserAlreadyExists) {
		return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
	}
	if err != nil {
		log.Printf("failed to signup: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to signup")
	}

	return c.NoContent(http.StatusCreated)
}

func (u *User) PostUsersSignin(c echo.Context) error {
	var userRequest openapi.PostUsersSigninJSONRequestBody
	err := c.Bind(&userRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json body")
	}

	name := values.NewUserName(string(userRequest.Name))
	err = name.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user name")
	}

	password := values.NewUserPassword([]byte(userRequest.Password))
	err = password.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user password")
	}

	user, err := u.authorizationService.Login(
		c.Request().Context(),
		name,
		password,
	)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user or password")
	}
	if err != nil {
		log.Printf("failed to signin: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to signin")
	}

	session, err := u.Session.getSession(c)
	if err != nil {
		log.Printf("failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	u.Session.setUser(session, user)

	err = u.Session.save(c, session)
	if err != nil {
		log.Printf("failed to save session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save session")
	}

	return c.NoContent(http.StatusOK)
}
