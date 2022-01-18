package domain

import (
	"time"

	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Response struct {
	id        values.ResponseID
	createdAt time.Time
}

func NewResponse(
	id values.ResponseID,
	createdAt time.Time,
) *Response {
	return &Response{
		id:        id,
		createdAt: createdAt,
	}
}

func (r *Response) GetID() values.ResponseID {
	return r.id
}

func (r *Response) GetCreatedAt() time.Time {
	return r.createdAt
}
