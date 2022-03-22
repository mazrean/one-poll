package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
	"github.com/mazrean/one-poll/service"
)

type Response struct {
	db                 repository.DB
	pollRepository     repository.Poll
	responseRepository repository.Response
	commentRepository  repository.Comment
	choiceRepository   repository.Choice
	pollAuthority      *PollAuthority
}

func NewResponse(
	db repository.DB,
	pollRepository repository.Poll,
	responseRepository repository.Response,
	commentRepository repository.Comment,
	choiceRepository repository.Choice,
	pollAuthority *PollAuthority,
) *Response {
	return &Response{
		db:                 db,
		pollRepository:     pollRepository,
		responseRepository: responseRepository,
		commentRepository:  commentRepository,
		choiceRepository:   choiceRepository,
		pollAuthority:      pollAuthority,
	}
}

func (r *Response) CreateResponse(
	ctx context.Context,
	user *domain.User,
	pollID values.PollID,
	choiceIDs []values.ChoiceID,
	commentContent values.CommentContent,
) (*service.ResponseInfo, error) {
	var (
		poll     *domain.Poll
		response *domain.Response
		comment  *domain.Comment
		choices  []*domain.Choice
	)
	err := r.db.Transaction(ctx, nil, func(ctx context.Context) error {
		pollInfo, err := r.pollRepository.GetPoll(ctx, pollID, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoPoll
		}
		if err != nil {
			return fmt.Errorf("failed to get poll: %w", err)
		}
		poll = pollInfo.Poll

		canResponse, err := r.pollAuthority.CanResponse(ctx, user, pollInfo.Owner, poll)
		if err != nil {
			return fmt.Errorf("failed to check can response: %w", err)
		}

		if !canResponse {
			return service.ErrResponseAlreadyExists
		}

		pollChoices, err := r.choiceRepository.GetChoicesByPollID(ctx, pollID, repository.LockTypeRecord)
		if err != nil {
			return fmt.Errorf("failed to get choices: %w", err)
		}

		pollChoiceMap := make(map[values.ChoiceID]*domain.Choice, len(pollChoices))
		for _, pollChoice := range pollChoices {
			pollChoiceMap[pollChoice.GetID()] = pollChoice
		}

		switch poll.GetPollType() {
		case values.PollTypeRadio:
			if len(choiceIDs) != 1 {
				return service.ErrTooManyChoice
			}
		}

		choices = make([]*domain.Choice, 0, len(choiceIDs))
		choiceMap := make(map[values.ChoiceID]struct{}, len(choiceIDs))
		for _, choiceID := range choiceIDs {
			choice, ok := pollChoiceMap[choiceID]
			if !ok {
				return service.ErrNoChoice
			}

			if _, ok := choiceMap[choiceID]; ok {
				return service.ErrDuplicateChoices
			}
			choiceMap[choiceID] = struct{}{}

			choices = append(choices, choice)
		}

		response = domain.NewResponse(
			values.NewResponseID(),
			time.Now(),
		)

		err = r.responseRepository.CreateResponse(ctx, user.GetID(), pollID, response, choiceIDs)
		if err != nil {
			return fmt.Errorf("failed to create response: %w", err)
		}

		if commentContent != "" {
			comment = domain.NewComment(
				values.NewCommentID(),
				commentContent,
			)

			err = r.commentRepository.CreateComment(ctx, response.GetID(), comment)
			if err != nil {
				return fmt.Errorf("failed to create comment: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return &service.ResponseInfo{
		Poll:     poll,
		Response: response,
		Comment:  comment,
		Choices:  choices,
	}, nil
}

func (r *Response) GetResult(ctx context.Context, user *domain.User, pollID values.PollID) (*service.Result, error) {
	pollInfo, err := r.pollRepository.GetPoll(ctx, pollID, repository.LockTypeNone)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrNoPoll
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	canResponse, err := r.pollAuthority.CanRead(ctx, user, pollInfo.Owner, pollInfo.Poll)
	if err != nil {
		return nil, fmt.Errorf("failed to check can response: %w", err)
	}

	if !canResponse {
		return nil, service.ErrForbidden
	}

	choices, err := r.choiceRepository.GetChoicesByPollID(ctx, pollID, repository.LockTypeNone)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}

	result, err := r.responseRepository.GetStatistics(ctx, pollID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	resultItems := make([]*service.ResultItem, 0, len(choices))
	for _, choice := range choices {
		resultItems = append(resultItems, &service.ResultItem{
			Choice: choice,
			Count:  result.ChoiceCount[choice.GetID()],
		})
	}

	return &service.Result{
		Poll:  pollInfo.Poll,
		Count: result.Count,
		Items: resultItems,
	}, nil
}
