package v1

import "github.com/cs-sysimpl/suzukake/service"

type Response struct {
	*Session
	responseService service.Response
}

func NewResponse(session *Session, responseService service.Response) *Response {
	return &Response{
		Session:         session,
		responseService: responseService,
	}
}
