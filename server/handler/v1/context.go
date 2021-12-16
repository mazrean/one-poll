package v1

import (
	"errors"
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	sessionContextKey = "session"
)

type Context struct {
	*Session
}

func NewContext(session *Session) *Context {
	return &Context{
		Session: session,
	}
}

//nolint:unused
func (ctx *Context) getSession(c echo.Context) (*sessions.Session, error) {
	iSession := c.Get(sessionContextKey)
	if iSession == nil {
		session, err := ctx.Session.getSession(c)
		if err != nil {
			return nil, fmt.Errorf("failed to get session: %w", err)
		}

		return session, nil
	}

	session, ok := iSession.(*sessions.Session)
	if !ok {
		return nil, errors.New("invalid session")
	}

	return session, nil
}
