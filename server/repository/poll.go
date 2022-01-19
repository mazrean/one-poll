package repository

import (
	"context"
	"database/sql"

	"github.com/cs-sysimpl/suzukake/domain"
	"github.com/cs-sysimpl/suzukake/domain/values"
)

type Poll interface {
	CreatePoll(ctx context.Context, poll *domain.Poll, owner values.UserID) error
	UpdatePollDeadline(ctx context.Context, id values.PollID, deadline sql.NullTime) error
	// GetPollsメソッドは Get polls,Get /users/me/owners Get/users/me/answers での仕様を想定
	// 現状のクライアントからの利用方法的に、一度にすべてのパラメタが埋まることはないが、仕様としてはパラメタ条件のAndをとる。
	GetPolls(ctx context.Context, params PollSearchParams) ([]*PollInfo, error)
	GetPoll(ctx context.Context, id values.PollID, lockType LockType) (*PollInfo, error)
	DeletePoll(ctx context.Context, id values.PollID) error
	AddTags(ctx context.Context, pollID values.PollID, tags []values.TagID) error
}

type PollInfo struct {
	*domain.Poll
	Owner *domain.User
}

// パラメタを指定していない要素については、その要素へのポインタをnilにする。
// 例: Matchに空文字が指定されていても、nilではないので、条件として働きます。
type PollSearchParams struct {
	Limit  *int
	Offset *int
	Match  *string
	Owner  *domain.User
	Answer *domain.User
}
