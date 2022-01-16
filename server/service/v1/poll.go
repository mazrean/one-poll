package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/repository"
	"github.com/cs-sysimpl/suzukake/service"
)

type Poll struct {
	db                 repository.DB
	pollRepository     repository.Poll
	choiceRepository   repository.Choice
	tagRepository      repository.Tag
	responseRepository repository.Response
}

func NewPoll(
	db repository.DB,
	pollRepository repository.Poll,
	choiceRepository repository.Choice,
	tagRepository repository.Tag,
	responseRepository repository.Response,
) *Poll {
	return &Poll{
		db:                 db,
		pollRepository:     pollRepository,
		choiceRepository:   choiceRepository,
		tagRepository:      tagRepository,
		responseRepository: responseRepository,
	}
}

func (p *Poll) CreatePoll(
	ctx context.Context,
	user *domain.User,
	title values.PollTitle,
	pollType values.PollType,
	deadline *time.Time,
	choiceLabels []values.ChoiceLabel,
	tagNames []values.TagName,
) (*service.PollInfo, error) {
	var (
		poll    *domain.Poll
		choices []*domain.Choice
		tags    []*domain.Tag
	)
	err := p.db.Transaction(ctx, nil, func(ctx context.Context) error {
		if deadline == nil {
			poll = domain.NewPollWithoutDeadLine(
				values.NewPollID(),
				title,
				pollType,
				time.Now(),
			)
		} else {
			poll = domain.NewPollWithDeadLine(
				values.NewPollID(),
				title,
				pollType,
				*deadline,
				time.Now(),
			)
		}

		err := p.pollRepository.CreatePoll(ctx, poll, user.GetID())
		if err != nil {
			return fmt.Errorf("failed to create poll: %w", err)
		}

		choices = make([]*domain.Choice, 0, len(choiceLabels))
		for _, choiceLabel := range choiceLabels {
			choices = append(choices, domain.NewChoice(
				values.NewChoiceID(),
				choiceLabel,
			))
		}

		err = p.choiceRepository.CreateChoices(ctx, poll.GetID(), choices)
		if err != nil {
			return fmt.Errorf("failed to create choices: %w", err)
		}

		tags, err = p.tagRepository.GetTagsByName(ctx, tagNames, repository.LockTypeRecord)
		if err != nil {
			return fmt.Errorf("failed to get tags: %w", err)
		}

		if len(tags) != len(tagNames) {
			return service.ErrNoTag
		}

		tagIDs := make([]values.TagID, 0, len(tags))
		for _, tag := range tags {
			tagIDs = append(tagIDs, tag.GetID())
		}

		err = p.pollRepository.AddTags(ctx, poll.GetID(), tagIDs)
		if err != nil {
			return fmt.Errorf("failed to add tags: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return &service.PollInfo{
		Poll:     poll,
		Choices:  choices,
		Tags:     tags,
		Owner:    user,
		Response: nil,
	}, nil
}

func (p *Poll) GetPolls(ctx context.Context, user *domain.User, params *service.PollSearchParams) ([]*service.PollInfo, error) {
	var repositoryParams *repository.PollSearchParams
	if params != nil {
		repositoryParams = &repository.PollSearchParams{
			Limit:  params.Limit,
			Offset: params.Offset,
			Match:  params.Match,
		}
	}

	polls, err := p.pollRepository.GetPolls(ctx, repositoryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get polls: %w", err)
	}

	pollIDs := make([]values.PollID, 0, len(polls))
	for _, poll := range polls {
		pollIDs = append(pollIDs, poll.Poll.GetID())
	}

	tagMap, err := p.tagRepository.GetTagsByPollIDs(ctx, pollIDs, repository.LockTypeNone)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	choiceMap, err := p.choiceRepository.GetChoicesByPollIDs(ctx, pollIDs, repository.LockTypeNone)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}

	responseMap, err := p.responseRepository.GetResponsesByUserIDAndPollIDs(ctx, user.GetID(), pollIDs, repository.LockTypeNone)
	if err != nil {
		return nil, fmt.Errorf("failed to get responses: %w", err)
	}

	pollInfos := make([]*service.PollInfo, 0, len(polls))
	for _, poll := range polls {
		choices, ok := choiceMap[poll.Poll.GetID()]
		if !ok {
			choices = []*domain.Choice{}
		}

		tags, ok := tagMap[poll.Poll.GetID()]
		if !ok {
			tags = []*domain.Tag{}
		}

		response, ok := responseMap[poll.Poll.GetID()]
		if !ok {
			response = nil
		}

		pollInfo := &service.PollInfo{
			Poll:     poll.Poll,
			Choices:  choices,
			Tags:     tags,
			Owner:    poll.Owner,
			Response: response,
		}
		pollInfos = append(pollInfos, pollInfo)
	}

	return pollInfos, nil
}

func (p *Poll) GetPoll(ctx context.Context, user *domain.User, id values.PollID) (*service.PollInfo, error) {
	poll, err := p.pollRepository.GetPoll(ctx, id, repository.LockTypeNone)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrNoPoll
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	tags, err := p.tagRepository.GetTagsByPollID(ctx, id, repository.LockTypeNone)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	choices, err := p.choiceRepository.GetChoicesByPollID(ctx, id, repository.LockTypeNone)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}

	response, err := p.responseRepository.GetResponseByUserIDAndPollID(ctx, user.GetID(), id, repository.LockTypeNone)
	if errors.Is(err, repository.ErrRecordNotFound) {
		response = nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	return &service.PollInfo{
		Poll:     poll.Poll,
		Choices:  choices,
		Tags:     tags,
		Owner:    poll.Owner,
		Response: response,
	}, nil
}

func (p *Poll) ClosePoll(ctx context.Context, user *domain.User, id values.PollID) error {
	err := p.db.Transaction(ctx, nil, func(ctx context.Context) error {
		poll, err := p.pollRepository.GetPoll(ctx, id, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoPoll
		}
		if err != nil {
			return fmt.Errorf("failed to get poll: %w", err)
		}

		if poll.Poll.IsExpired() {
			return service.ErrPollClosed
		}

		if poll.Owner.GetID() != user.GetID() {
			return service.ErrNotOwner
		}

		err = p.pollRepository.UpdatePollDeadline(ctx, id, sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		})
		if errors.Is(err, repository.ErrNoRecordUpdated) {
			return service.ErrNoPoll
		}
		if err != nil {
			return fmt.Errorf("failed to close poll: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}

func (p *Poll) DeletePoll(ctx context.Context, user *domain.User, id values.PollID) error {
	err := p.pollRepository.DeletePoll(ctx, id)
	if errors.Is(err, repository.ErrNoRecordDeleted) {
		return service.ErrNoPoll
	}
	if err != nil {
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	return nil
}
