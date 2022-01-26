package gorm2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/repository"
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

func (p *Poll) UpdatePollDeadline(ctx context.Context, id values.PollID, deadline sql.NullTime) error {
	db, err := p.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	result := db.
		Model(&PollTable{}).
		Where("id = ?", uuid.UUID(id)).
		Update("deadline", deadline)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update poll deadline: %w", err)
	}

	if result.RowsAffected == 0 {
		return repository.ErrNoRecordUpdated
	}

	return nil
}

func (p *Poll) GetPolls(ctx context.Context, params *repository.PollSearchParams) ([]*repository.PollInfo, error) {
	db, err := p.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var pollTables []PollTable
	query := db.
		Joins("Owner").
		Joins("PollType").
		Order("created_at DESC")

	if params != nil {
		if params.Match != "" {
			query = query.Where("polls.title like ?", "%"+params.Match+"%")
		}

		if params.Limit > 0 {
			query = query.Limit(params.Limit)
		} else if params.Limit < 0 {
			return nil, repository.ErrInvalidParameterValue("Limit", "be positive")
		}

		if params.Offset > 0 {
			query = query.Offset(params.Offset)
		} else if params.Offset < 0 {
			return nil, repository.ErrInvalidParameterValue("Offset", "be positive")
		}

		if params.Owner != nil {
			id := uuid.UUID(params.Owner.GetID()).String()
			query = query.Where("polls.owner_id = ?", id)
		}

		if params.Answer != nil {
			var responseTable []ResponseTable
			id := uuid.UUID(params.Answer.GetID())
			err = db.Where("respondent_id", id).Select("poll_id").Find(&responseTable).Error
			if err != nil {
				return nil, fmt.Errorf("failed to get responses: %w", err)
			}
			uuidPollIDs := make([]uuid.UUID, 0, len(responseTable))
			for _, responseTable := range responseTable {
				uuidPollIDs = append(uuidPollIDs, responseTable.PollID)
			}
			query = query.Where("polls.id IN ?", uuidPollIDs)
		}
	}

	err = query.Find(&pollTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get polls: %w", err)
	}

	pollInfos := make([]*repository.PollInfo, 0, len(pollTables))
	for _, pollTable := range pollTables {
		var polltype values.PollType
		switch pollTable.PollType.Name {
		case pollTypeRadio:
			polltype = values.PollTypeRadio
		default:
			return nil, fmt.Errorf("unsupported poll type: %s", pollTable.PollType.Name)
		}

		var poll *domain.Poll
		if pollTable.Deadline.Valid {
			poll = domain.NewPollWithDeadLine(
				values.NewPollIDFromUUID(pollTable.ID),
				values.NewPollTitle(pollTable.Title),
				polltype,
				pollTable.Deadline.Time,
				pollTable.CreatedAt,
			)
		} else {
			poll = domain.NewPollWithoutDeadLine(
				values.NewPollIDFromUUID(pollTable.ID),
				values.NewPollTitle(pollTable.Title),
				polltype,
				pollTable.CreatedAt,
			)
		}

		pollInfos = append(pollInfos, &repository.PollInfo{
			Poll: poll,
			Owner: domain.NewUser(
				values.NewUserIDFromUUID(pollTable.Owner.ID),
				values.NewUserName(pollTable.Owner.Name),
				values.NewUserHashedPassword([]byte(pollTable.Owner.Password)),
			),
		})
	}

	return pollInfos, nil
}

func (p *Poll) GetPoll(ctx context.Context, id values.PollID, lockType repository.LockType) (*repository.PollInfo, error) {
	db, err := p.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = p.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var pollTable PollTable
	query := db.
		Joins("Owner").
		Joins("PollType").
		Where("polls.id = ?", uuid.UUID(id))

	err = query.Take(&pollTable).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	var polltype values.PollType
	switch pollTable.PollType.Name {
	case pollTypeRadio:
		polltype = values.PollTypeRadio
	default:
		return nil, fmt.Errorf("unsupported poll type: %s", pollTable.PollType.Name)
	}

	var poll *domain.Poll
	if pollTable.Deadline.Valid {
		poll = domain.NewPollWithDeadLine(
			values.NewPollIDFromUUID(pollTable.ID),
			values.NewPollTitle(pollTable.Title),
			polltype,
			pollTable.Deadline.Time,
			pollTable.CreatedAt,
		)
	} else {
		poll = domain.NewPollWithoutDeadLine(
			values.NewPollIDFromUUID(pollTable.ID),
			values.NewPollTitle(pollTable.Title),
			polltype,
			pollTable.CreatedAt,
		)
	}

	return &repository.PollInfo{
		Poll: poll,
		Owner: domain.NewUser(
			values.NewUserIDFromUUID(pollTable.Owner.ID),
			values.NewUserName(pollTable.Owner.Name),
			values.NewUserHashedPassword([]byte(pollTable.Owner.Password)),
		),
	}, nil
}

func (p *Poll) DeletePoll(ctx context.Context, id values.PollID) error {
	db, err := p.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	result := db.
		Where("id = ?", uuid.UUID(id)).
		Delete(&PollTable{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	if result.RowsAffected == 0 {
		return repository.ErrNoRecordDeleted
	}

	return nil
}

func (p *Poll) AddTags(ctx context.Context, pollID values.PollID, tagIDs []values.TagID) error {
	db, err := p.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	tagTables := make([]TagTable, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		tagTables = append(tagTables, TagTable{
			ID: uuid.UUID(tagID),
		})
	}

	err = db.
		Model(&PollTable{
			ID: uuid.UUID(pollID),
		}).
		Omit("Tags.*").
		Association("Tags").
		Append(tagTables)
	if err != nil {
		return fmt.Errorf("failed to add tags: %w", err)
	}

	return nil
}
