package repository

import (
	"context"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
)

type WebauthnCredential interface {
	StoreCredential(ctx context.Context, userID values.UserID, credential *domain.WebAuthnCredential) error
	GetCredentialsByUserID(ctx context.Context, userID values.UserID) ([]*domain.WebAuthnCredential, error)
	GetCredentialWithUserByCredID(ctx context.Context, credID values.WebAuthnCredentialCredID) (*domain.WebAuthnCredential, *domain.User, error)
}
