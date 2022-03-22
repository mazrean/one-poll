package service

import (
	"context"

	"github.com/mazrean/one-poll/domain"
)

type Tag interface {
	GetTags(ctx context.Context) ([]*domain.Tag, error)
}
