package domain

import (
	"time"

	"github.com/mazrean/one-poll/domain/values"
)

type WebauthnCredential struct {
	id         values.WebauthnCredentialID
	credID     values.WebauthnCredentialCredID
	name       values.WebauthnCredentialName
	publicKey  values.WebauthnCredentialPublicKey
	algorithm  values.WebauthnCredentialAlgorithm
	transports []values.WebauthnCredentialTransport
	createdAt  time.Time
}

func NewWebauthnCredential(
	id values.WebauthnCredentialID,
	credID values.WebauthnCredentialCredID,
	name values.WebauthnCredentialName,
	publicKey values.WebauthnCredentialPublicKey,
	algorithm values.WebauthnCredentialAlgorithm,
	transports []values.WebauthnCredentialTransport,
	createdAt time.Time,
) *WebauthnCredential {
	return &WebauthnCredential{
		id:         id,
		credID:     credID,
		name:       name,
		publicKey:  publicKey,
		algorithm:  algorithm,
		transports: transports,
		createdAt:  createdAt,
	}
}

func (w *WebauthnCredential) ID() values.WebauthnCredentialID {
	return w.id
}

func (w *WebauthnCredential) CredID() values.WebauthnCredentialCredID {
	return w.credID
}

func (w *WebauthnCredential) Name() values.WebauthnCredentialName {
	return w.name
}

func (w *WebauthnCredential) PublicKey() values.WebauthnCredentialPublicKey {
	return w.publicKey
}

func (w *WebauthnCredential) Algorithm() values.WebauthnCredentialAlgorithm {
	return w.algorithm
}

func (w *WebauthnCredential) Transports() []values.WebauthnCredentialTransport {
	return w.transports
}

func (w *WebauthnCredential) CreatedAt() time.Time {
	return w.createdAt
}
