package domain

import "github.com/cs-sysimpl/suzukake/domain/values"

type Comment struct {
	id      values.CommentID
	content values.CommentContent
}

func NewComment(
	id values.CommentID,
	content values.CommentContent,
) Comment {
	return Comment{
		id:      id,
		content: content,
	}
}

func (c *Comment) GetID() values.CommentID {
	return c.id
}

func (c *Comment) GetContent() values.CommentContent {
	return c.content
}
