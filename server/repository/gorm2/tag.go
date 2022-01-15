package gorm2

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/repository"
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

func (t *Tag) GetTagsByName(ctx context.Context, names []values.TagName, lockType repository.LockType) ([]*domain.Tag, error) {
	db, err := t.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = t.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	strNames := make([]string, 0, len(names))
	for _, name := range names {
		strNames = append(strNames, string(name))
	}

	var tagTables []TagTable
	err = db.
		Where("name IN ?", strNames).
		Select("id", "name").
		Find(&tagTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	tags := make([]*domain.Tag, 0, len(tagTables))
	for _, tagTable := range tagTables {
		tags = append(tags, domain.NewTag(
			values.NewTagIDFromUUID(tagTable.ID),
			values.NewTagName(tagTable.Name),
		))
	}

	return tags, nil
}
