package gorm2

import (
	"context"
	"errors"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Response struct {
	db *DB
}

func NewResponse(db *DB) *Response {
	return &Response{
		db: db,
	}
}

func (r *Response) CreateResponse(ctx context.Context, userID values.UserID, pollID values.PollID, response *domain.Response, choiceIDs []values.ChoiceID) error {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	responseTable := ResponseTable{
		ID:           uuid.UUID(response.GetID()),
		PollID:       uuid.UUID(pollID),
		RespondentID: uuid.UUID(userID),
		CreatedAt:    response.GetCreatedAt(),
	}

	err = db.Create(&responseTable).Error
	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	choiceTables := make([]ChoiceTable, 0, len(choiceIDs))
	for _, choiceID := range choiceIDs {
		choiceTables = append(choiceTables, ChoiceTable{
			ID: uuid.UUID(choiceID),
		})
	}

	err = db.
		Model(&responseTable).
		Association("Choices").
		Append(choiceTables)
	if err != nil {
		return fmt.Errorf("failed to append choices: %w", err)
	}

	return nil
}

func (r *Response) GetResponsesByPollID(ctx context.Context, pollID values.PollID) ([]*repository.ResponseInfo, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var responseTables []ResponseTable
	err = db.
		Where("poll_id = ?", uuid.UUID(pollID)).
		Joins("Respondent").
		Order("created_at DESC").
		Select("responses.id", "created_at").
		Find(&responseTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}

	responseInfos := make([]*repository.ResponseInfo, 0, len(responseTables))
	for _, responseTable := range responseTables {
		responseInfos = append(responseInfos,
			&repository.ResponseInfo{
				Response: domain.NewResponse(
					values.NewResponseIDFromUUID(responseTable.ID),
					responseTable.CreatedAt,
				),
				Respondent: domain.NewUser(
					values.NewUserIDFromUUID(responseTable.Respondent.ID),
					values.NewUserName(responseTable.Respondent.Name),
					values.NewUserHashedPassword([]byte(responseTable.Respondent.Password)),
				),
			},
		)
	}

	return responseInfos, nil
}

func (r *Response) GetResponseByUserIDAndPollID(ctx context.Context, userID values.UserID, pollID values.PollID, lockType repository.LockType) (*domain.Response, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = r.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var responseTable ResponseTable
	err = db.
		Where("poll_id = ? AND respondent_id = ?", uuid.UUID(pollID), uuid.UUID(userID)).
		Select("id", "created_at").
		Take(&responseTable).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	return domain.NewResponse(
		values.NewResponseIDFromUUID(responseTable.ID),
		responseTable.CreatedAt,
	), nil
}

func (r *Response) GetResponsesByUserIDAndPollIDs(ctx context.Context, userID values.UserID, pollIDs []values.PollID, lockType repository.LockType) (map[values.PollID]*domain.Response, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = r.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	uuidPollIDs := make([]uuid.UUID, 0, len(pollIDs))
	for _, pollID := range pollIDs {
		uuidPollIDs = append(uuidPollIDs, uuid.UUID(pollID))
	}

	var responseTables []ResponseTable
	err = db.
		Where("poll_id IN ? AND respondent_id = ?", uuidPollIDs, uuid.UUID(userID)).
		Select("id", "poll_id", "created_at").
		Find(&responseTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}

	responses := make(map[values.PollID]*domain.Response, len(responseTables))
	for _, responseTable := range responseTables {
		responses[values.NewPollIDFromUUID(responseTable.PollID)] = domain.NewResponse(
			values.NewResponseIDFromUUID(responseTable.ID),
			responseTable.CreatedAt,
		)
	}

	return responses, nil
}

func (r *Response) GetStatistics(ctx context.Context, pollID values.PollID) (*repository.Statistics, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var responseCount int64
	err = db.
		Session(&gorm.Session{}).
		Model(&ResponseTable{}).
		Where("poll_id = ?", uuid.UUID(pollID)).
		Count(&responseCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}

	var choiceCounts []struct {
		ChoiceID uuid.UUID
		Count    int64
	}
	err = db.
		Model(&ResponseTable{}).
		Where("responses.poll_id = ?", uuid.UUID(pollID)).
		Joins("INNER JOIN response_choice_relations ON responses.id = response_choice_relations.response_table_id").
		Joins("INNER JOIN choices ON response_choice_relations.choice_table_id = choices.id").
		Group("choices.id").
		Select("choices.id AS choice_id, COUNT(responses.id) AS count").
		Find(&choiceCounts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get choice counts: %w", err)
	}

	choiceCountMap := make(map[values.ChoiceID]int, len(choiceCounts))
	for _, choiceCount := range choiceCounts {
		choiceCountMap[values.NewChoiceIDFromUUID(choiceCount.ChoiceID)] = int(choiceCount.Count)
	}

	return &repository.Statistics{
		Count:       int(responseCount),
		ChoiceCount: choiceCountMap,
	}, nil
}
