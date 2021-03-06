package domain

import (
	"database/sql"
	"time"

	"github.com/mazrean/one-poll/domain/values"
)

type Poll struct {
	id        values.PollID
	title     values.PollTitle
	pollType  values.PollType
	deadline  sql.NullTime
	createdAt time.Time
}

func NewPollWithDeadLine(
	id values.PollID,
	title values.PollTitle,
	pollType values.PollType,
	deadline time.Time,
	createdAt time.Time,
) *Poll {
	return &Poll{
		id:       id,
		title:    title,
		pollType: pollType,
		deadline: sql.NullTime{
			Time:  deadline,
			Valid: true,
		},
		createdAt: createdAt,
	}
}

func NewPollWithoutDeadLine(
	id values.PollID,
	title values.PollTitle,
	pollType values.PollType,
	createdAt time.Time,
) *Poll {
	return &Poll{
		id:        id,
		title:     title,
		pollType:  pollType,
		createdAt: createdAt,
	}
}

func (p *Poll) GetID() values.PollID {
	return p.id
}

func (p *Poll) GetTitle() values.PollTitle {
	return p.title
}

func (p *Poll) GetPollType() values.PollType {
	return p.pollType
}

func (p *Poll) GetDeadline() sql.NullTime {
	return p.deadline
}

func (p *Poll) GetCreatedAt() time.Time {
	return p.createdAt
}

func (p *Poll) IsExpired() bool {
	return p.deadline.Valid && p.deadline.Time.Before(time.Now())
}
