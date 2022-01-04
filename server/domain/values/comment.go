package values

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type (
	CommentID      uuid.UUID
	CommentContent string
)

func NewCommentID() CommentID {
	return CommentID(uuid.New())
}

func NewCommentIDFromUUID(id uuid.UUID) CommentID {
	return CommentID(id)
}

func NewCommentContent(content string) CommentContent {
	return CommentContent(content)
}

var (
	ErrCommentContentTooLong = errors.New("comment content is too long")
)

func (tn CommentContent) Validate() error {
	if utf8.RuneCountInString(string(tn)) > 50 {
		return ErrCommentContentTooLong
	}

	return nil
}
