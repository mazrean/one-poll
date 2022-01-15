package gorm2

import (
	"context"
	"fmt"

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
