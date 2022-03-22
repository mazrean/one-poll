package domain

import "github.com/mazrean/one-poll/domain/values"

type Tag struct {
	id   values.TagID
	name values.TagName
}

func NewTag(
	id values.TagID,
	name values.TagName,
) *Tag {
	return &Tag{
		id:   id,
		name: name,
	}
}

func (t *Tag) GetID() values.TagID {
	return t.id
}

func (t *Tag) GetName() values.TagName {
	return t.name
}
