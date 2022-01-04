package v1

import (
	"github.com/cs-sysimpl/suzukake/service"
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
