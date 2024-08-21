package repository

import (
	"context"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
)

type WebAuthnCredential interface {
	StoreCredential(ctx context.Context, userID values.UserID, credential *domain.WebAuthnCredential) error
	GetCredentialsByUserID(ctx context.Context, userID values.UserID) ([]*domain.WebAuthnCredential, error)
	GetCredentialWithUserByCredID(ctx context.Context, credID values.WebAuthnCredentialCredID, lockType LockType) (*domain.WebAuthnCredential, *domain.User, error)
	UpdateLastUsedAt(ctx context.Context, credential *domain.WebAuthnCredential) error
	DeleteCredential(ctx context.Context, userID values.UserID, credID values.WebAuthnCredentialCredID) error
}
