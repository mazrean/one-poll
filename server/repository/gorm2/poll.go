package gorm2

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	pollTypeRadio = "radio"
)

type Poll struct {
	db        *DB
	pollTypes []PollTypeTable
}

func NewPoll(db *DB) (*Poll, error) {
	gormDB, err := db.getDB(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	pollTypes, err := setupPollTypeTable(gormDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup poll type table: %w", err)
	}

	return &Poll{
		db:        db,
		pollTypes: pollTypes,
	}, nil
}

func setupPollTypeTable(db *gorm.DB) ([]PollTypeTable, error) {
	pollTypes := []PollTypeTable{
		{
			Name:   pollTypeRadio,
			Active: true,
		},
	}

	for i, pollType := range pollTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", pollType.Name).
			FirstOrCreate(&pollType).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create resource type: %w", err)
		}

		pollTypes[i] = pollType
	}

	return pollTypes, nil
}

func (p *Poll) CreatePoll(ctx context.Context, poll *domain.Poll, ownerID values.UserID) error {
	db, err := p.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	var pollTypeID int
	for _, pollType := range p.pollTypes {
		switch poll.GetPollType() {
		case values.PollTypeRadio:
			if pollType.Name == pollTypeRadio {
				pollTypeID = pollType.ID
			}
		default:
			return fmt.Errorf("unsupported poll type: %d", poll.GetPollType())
		}
	}

	pollTable := PollTable{
		ID:        uuid.UUID(poll.GetID()),
		OwnerID:   uuid.UUID(ownerID),
		Title:     string(poll.GetTitle()),
		TypeID:    pollTypeID,
		Deadline:  poll.GetDeadline(),
		CreatedAt: poll.GetCreatedAt(),
	}
	err = db.Create(&pollTable).Error
	if err != nil {
		return fmt.Errorf("failed to create poll: %w", err)
	}

	return nil
}
