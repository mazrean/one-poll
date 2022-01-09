package v1

import (
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/pkg/common"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Session struct {
	key    string
	secret string
	store  sessions.Store
}

func NewSession(key common.SessionKey, secret common.SessionSecret) *Session {
	store := sessions.NewCookieStore([]byte(secret))

	/*
		gorilla/sessionsの内部で使われているgobが
		time.Timeのエンコード・デコードをできるようにRegisterする
	*/
	gob.Register(time.Time{})
	gob.Register(uuid.UUID{})

	return &Session{
		key:    string(key),
		secret: string(secret),
		store:  store,
	}
}

func (s *Session) Use(e *echo.Echo) {
	e.Use(session.Middleware(s.store))
}

//nolint:unused
func (s *Session) getSession(c echo.Context) (*sessions.Session, error) {
	session, err := s.store.Get(c.Request(), s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	c.Set(sessionContextKey, session)

	return session, nil
}

//nolint:unused
func (s *Session) save(c echo.Context, session *sessions.Session) error {
	err := s.store.Save(c.Request(), c.Response(), session)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

//nolint:unused
func (s *Session) revoke(session *sessions.Session) {
	session.Options.MaxAge = -1
}

var (
	ErrNoValue     = errors.New("no value")
	ErrValueBroken = errors.New("value broken")
)

const (
	userIDSessionKey             = "userID"
	userNameSessionKey           = "userName"
	userHashedPasswordSessionKey = "userHashedPassword"
)

func (s *Session) setUser(session *sessions.Session, user *domain.User) {
	session.Values[userIDSessionKey] = uuid.UUID(user.GetID())
	session.Values[userNameSessionKey] = string(user.GetName())
	session.Values[userHashedPasswordSessionKey] = []byte(user.GetHashedPassword())
}

func (s *Session) getUser(session *sessions.Session) (*domain.User, error) {
	iUserID, ok := session.Values[userIDSessionKey]
	if !ok {
		return nil, ErrNoValue
	}

	userID, ok := iUserID.(uuid.UUID)
	if !ok {
		return nil, ErrValueBroken
	}

	iUserName, ok := session.Values[userNameSessionKey]
	if !ok {
		return nil, ErrNoValue
	}

	userName, ok := iUserName.(string)
	if !ok {
		return nil, ErrValueBroken
	}

	iUserHashedPassword, ok := session.Values[userHashedPasswordSessionKey]
	if !ok {
		return nil, ErrNoValue
	}

	userHashedPassword, ok := iUserHashedPassword.([]byte)
	if !ok {
		return nil, ErrValueBroken
	}

	return domain.NewUser(
		values.NewUserIDFromUUID(userID),
		values.NewUserName(userName),
		values.NewUserHashedPassword(userHashedPassword),
	), nil
}
