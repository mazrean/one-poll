package gorm2

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/google/uuid"
)

type Response struct {
	db *DB
}

func NewResponse(db *DB) *Response {
	return &Response{
		db: db,
	}
}

func (r *Response) GetResponsesByPollID(ctx context.Context, pollID values.PollID) ([]*domain.Response, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var responseTables []ResponseTable
	err = db.
		Where("poll_id = ?", uuid.UUID(pollID)).
		Order("created_at DESC").
		Select("id", "created_at").
		Find(&responseTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}

	responses := make([]*domain.Response, 0, len(responseTables))
	for _, responseTable := range responseTables {
		responses = append(responses, domain.NewResponse(
			values.NewResponseIDFromUUID(responseTable.ID),
			responseTable.CreatedAt,
		))
	}

	return responses, nil
}
