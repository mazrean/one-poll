package service

import (
	"context"

	"github.com/cs-sysimpl/suzukake/domain"
)

type Tag interface {
	GetTags(ctx context.Context) ([]*domain.Tag, error)
}
