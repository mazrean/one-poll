package v1

import (
	"context"
	"fmt"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
	"github.com/cs-sysimpl/suzukake/repository"
	"github.com/cs-sysimpl/suzukake/service"
)

func (p *Poll) GetOwnerPolls(ctx context.Context, owner *domain.User) ([]*service.PollInfo, error) {
	if owner == nil {
		return nil, fmt.Errorf("owner parameter is required")
	}

	repositoryParams := &repository.PollSearchParams{
		Owner: owner,
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

		pollInfo := &service.PollInfo{
			Poll:     poll.Poll,
			Choices:  choices,
			Tags:     tags,
			Owner:    poll.Owner,
			Response: nil, // Ownerは自身の投票に対して回答できないため。
		}
		pollInfos = append(pollInfos, pollInfo)
	}
	return pollInfos, nil
}

func (p *Poll) GetAnsweredPolls(ctx context.Context, respondent *domain.User) ([]*service.PollInfo, error) {
	if respondent == nil {
		return nil, fmt.Errorf("respondent parameter is required")
	}

	repositoryParams := &repository.PollSearchParams{
		Answer: respondent,
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

	var responseMap map[values.PollID]*domain.Response

	// 答えたPollID(かつ自分の答えたPoll)のレスポンスを取得
	responseMap, err = p.responseRepository.GetResponsesByUserIDAndPollIDs(ctx, respondent.GetID(), pollIDs, repository.LockTypeNone)
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
