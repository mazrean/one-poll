package v1

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/repository"
)

type Tag struct {
	db            repository.DB
	tagRepository repository.Tag
}

func NewTag(db repository.DB, tagRepository repository.Tag) *Tag {
	return &Tag{
		db:            db,
		tagRepository: tagRepository,
	}
}

func (t *Tag) GetTags(ctx context.Context) ([]*domain.Tag, error) {
	tags, err := t.tagRepository.GetTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}
