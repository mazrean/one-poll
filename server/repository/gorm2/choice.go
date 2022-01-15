package gorm2

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/google/uuid"
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
	for _, choice := range choices {
		choiceTables = append(choiceTables, ChoiceTable{
			ID:     uuid.UUID(choice.GetID()),
			PollID: uuid.UUID(pollID),
			Name:   string(choice.GetLabel()),
		})
	}

	err = db.Create(&choiceTables).Error
	if err != nil {
		return fmt.Errorf("failed to create choices: %w", err)
	}

	return nil
}
