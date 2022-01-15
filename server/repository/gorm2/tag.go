package gorm2

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/google/uuid"
)

type Tag struct {
	db *DB
}

func NewTag(db *DB) *Tag {
	return &Tag{
		db: db,
	}
}

func (t *Tag) CreateTag(ctx context.Context, tag *domain.Tag) error {
	db, err := t.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	tagTable := TagTable{
		ID:   uuid.UUID(tag.GetID()),
		Name: string(tag.GetName()),
	}
	err = db.Create(&tagTable).Error
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}
