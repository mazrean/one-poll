package v1

import (
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/cs-sysimpl/suzukake/pkg/common"
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
