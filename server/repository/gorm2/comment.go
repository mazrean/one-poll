package gorm2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
)

type Comment struct {
	db *DB
}

func NewComment(db *DB) *Comment {
	return &Comment{
		db: db,
	}
}

func (c *Comment) CreateComment(ctx context.Context, responseID values.ResponseID, comment *domain.Comment) error {
	db, err := c.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	commentTable := CommentTable{
		ID:         uuid.UUID(comment.GetID()),
		ResponseID: uuid.UUID(responseID),
		Comment:    string(comment.GetContent()),
	}

	err = db.Create(&commentTable).Error
	if err != nil {
		return fmt.Errorf("failed to create poll: %w", err)
	}

	return nil
}

func (c *Comment) GetCommentsByResponseIDs(ctx context.Context, responseIDs []values.ResponseID, params repository.CommentGetParams) (map[values.ResponseID]*domain.Comment, error) {
	db, err := c.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	uuidResponseIDs := make([]uuid.UUID, 0, len(responseIDs))
	for _, responseID := range responseIDs {
		uuidResponseIDs = append(uuidResponseIDs, uuid.UUID(responseID))
	}

	var commentTables []CommentTable
	query := db.
		Where("response_id IN ?", uuidResponseIDs)

	if params.Limit != nil {
		if *params.Limit > 0 {
			query = query.Limit(*params.Limit)
		}
	}
	if params.Offset != nil {
		if *params.Offset > 0 {
			query = query.Offset(*params.Offset)
		}
	}

	err = query.Select("id", "response_id", "comment").Find(&commentTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	comments := make(map[values.ResponseID]*domain.Comment, len(commentTables))
	for _, commentTable := range commentTables {
		comments[values.NewResponseIDFromUUID(commentTable.ResponseID)] = domain.NewComment(
			values.NewCommentIDFromUUID(commentTable.ID),
			values.CommentContent(commentTable.Comment),
		)
	}

	return comments, nil
}
