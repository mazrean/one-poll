package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mazrean/one-poll/domain/values"
	openapi "github.com/mazrean/one-poll/handler/v1/openapi"
	"github.com/mazrean/one-poll/service"
	"github.com/oapi-codegen/runtime/types"
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

	user, err := u.authorizationService.Signup(
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

func (u *User) PostUsersSignout(c echo.Context) error {
	session, err := u.Session.getSession(c)
	if err != nil {
		log.Printf("failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	_, err = u.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	u.revoke(session)

	err = u.Session.save(c, session)
	if err != nil {
		log.Printf("failed to save session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save session")
	}

	return c.NoContent(http.StatusOK)
}

func (u *User) GetUsersMe(c echo.Context) error {
	session, err := u.Session.getSession(c)
	if err != nil {
		log.Printf("failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	user, err := u.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	return c.JSON(http.StatusOK, openapi.User{
		Name: openapi.UserName(user.GetName()),
		Uuid: types.UUID(user.GetID()),
	})
}

func (u *User) DeleteUsersMe(c echo.Context) error {
	session, err := u.Session.getSession(c)
	if err != nil {
		log.Printf("failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	user, err := u.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	err = u.authorizationService.DeleteAccount(
		c.Request().Context(),
		user,
	)
	if err != nil {
		log.Printf("failed to delete user: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user")
	}

	u.revoke(session)

	err = u.Session.save(c, session)
	if err != nil {
		log.Printf("failed to save session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save session")
	}

	return c.NoContent(http.StatusOK)
}
