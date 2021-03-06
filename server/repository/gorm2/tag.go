package gorm2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
	"gorm.io/gorm"
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

func (t *Tag) GetTags(ctx context.Context) ([]*domain.Tag, error) {
	db, err := t.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var tagTables []TagTable
	err = db.
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

func (t *Tag) GetTagsByPollIDs(ctx context.Context, pollIDs []values.PollID, lockType repository.LockType) (map[values.PollID][]*domain.Tag, error) {
	db, err := t.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = t.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	uuidPollIDs := make([]uuid.UUID, 0, len(pollIDs))
	for _, pollID := range pollIDs {
		uuidPollIDs = append(uuidPollIDs, uuid.UUID(pollID))
	}

	var pollTables []PollTable
	err = db.
		Where("polls.id IN ?", uuidPollIDs).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Select("polls.id").
		Find(&pollTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	tags := make(map[values.PollID][]*domain.Tag, len(pollTables))
	for _, pollTable := range pollTables {
		pollID := values.NewPollIDFromUUID(pollTable.ID)
		tags[pollID] = make([]*domain.Tag, 0, len(pollTable.Tags))
		for _, tag := range pollTable.Tags {
			tags[pollID] = append(tags[pollID], domain.NewTag(
				values.NewTagIDFromUUID(tag.ID),
				values.NewTagName(tag.Name),
			))
		}
	}

	return tags, nil
}

func (t *Tag) GetTagsByPollID(ctx context.Context, pollID values.PollID, lockType repository.LockType) ([]*domain.Tag, error) {
	db, err := t.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = t.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var pollTable PollTable
	err = db.
		Where("polls.id = ?", uuid.UUID(pollID)).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Select("polls.id").
		Take(&pollTable).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	tags := make([]*domain.Tag, 0, len(pollTable.Tags))
	for _, tag := range pollTable.Tags {
		tags = append(tags, domain.NewTag(
			values.NewTagIDFromUUID(tag.ID),
			values.NewTagName(tag.Name),
		))
	}

	return tags, nil
}
