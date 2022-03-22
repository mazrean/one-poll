package gorm2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
)

type Choice struct {
	db *DB
}

func NewChoice(db *DB) *Choice {
	return &Choice{
		db: db,
	}
}

func (c *Choice) CreateChoices(ctx context.Context, pollID values.PollID, choices []*domain.Choice) error {
	db, err := c.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	choiceTables := make([]ChoiceTable, 0, len(choices))
	for i, choice := range choices {
		choiceTables = append(choiceTables, ChoiceTable{
			ID:     uuid.UUID(choice.GetID()),
			PollID: uuid.UUID(pollID),
			Name:   string(choice.GetLabel()),
			Order:  uint8(i),
		})
	}

	err = db.Create(&choiceTables).Error
	if err != nil {
		return fmt.Errorf("failed to create choices: %w", err)
	}

	return nil
}

func (c *Choice) GetChoices(ctx context.Context, ids []values.ChoiceID, lockType repository.LockType) ([]*domain.Choice, error) {
	db, err := c.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	uuidChoiceIDs := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uuidChoiceIDs = append(uuidChoiceIDs, uuid.UUID(id))
	}

	var choiceTables []ChoiceTable
	err = db.
		Where("id IN ?", uuidChoiceIDs).
		Order("`order`").
		Select("id", "name").
		Find(&choiceTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}

	choices := make([]*domain.Choice, 0, len(choiceTables))
	for _, choiceTable := range choiceTables {
		choices = append(choices, domain.NewChoice(
			values.NewChoiceIDFromUUID(choiceTable.ID),
			values.NewChoiceLabel(choiceTable.Name),
		))
	}

	return choices, nil
}

func (c *Choice) GetChoicesByPollIDs(ctx context.Context, pollIDs []values.PollID, lockType repository.LockType) (map[values.PollID][]*domain.Choice, error) {
	db, err := c.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = c.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	uuidPollIDs := make([]uuid.UUID, 0, len(pollIDs))
	for _, id := range pollIDs {
		uuidPollIDs = append(uuidPollIDs, uuid.UUID(id))
	}

	var choiceTables []ChoiceTable
	err = db.
		Where("poll_id IN ?", uuidPollIDs).
		Order("`order`").
		Select("id", "name", "poll_id").
		Find(&choiceTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}

	choices := make(map[values.PollID][]*domain.Choice, len(choiceTables))
	for _, choiceTable := range choiceTables {
		choices[values.NewPollIDFromUUID(choiceTable.PollID)] = append(choices[values.NewPollIDFromUUID(choiceTable.PollID)], domain.NewChoice(
			values.NewChoiceIDFromUUID(choiceTable.ID),
			values.NewChoiceLabel(choiceTable.Name),
		))
	}

	return choices, nil
}

func (c *Choice) GetChoicesByPollID(ctx context.Context, pollID values.PollID, lockType repository.LockType) ([]*domain.Choice, error) {
	db, err := c.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = c.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var choiceTables []ChoiceTable
	err = db.
		Where("poll_id = ?", uuid.UUID(pollID)).
		Order("`order`").
		Select("id", "name").
		Find(&choiceTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}

	choices := make([]*domain.Choice, 0, len(choiceTables))
	for _, choiceTable := range choiceTables {
		choices = append(choices, domain.NewChoice(
			values.NewChoiceIDFromUUID(choiceTable.ID),
			values.NewChoiceLabel(choiceTable.Name),
		))
	}

	return choices, nil
}
