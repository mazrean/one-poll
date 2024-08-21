package domain

import (
	"time"

	"github.com/mazrean/one-poll/domain/values"
)

type WebAuthnCredential struct {
	id         values.WebAuthnCredentialID
	credID     values.WebAuthnCredentialCredID
	name       values.WebAuthnCredentialName
	publicKey  values.WebAuthnCredentialPublicKey
	algorithm  values.WebAuthnCredentialAlgorithm
	transports []values.WebAuthnCredentialTransport
	createdAt  time.Time
	lastUsedAt time.Time
}

func NewWebAuthnCredential(
	id values.WebAuthnCredentialID,
	credID values.WebAuthnCredentialCredID,
	name values.WebAuthnCredentialName,
	publicKey values.WebAuthnCredentialPublicKey,
	algorithm values.WebAuthnCredentialAlgorithm,
	transports []values.WebAuthnCredentialTransport,
	createdAt time.Time,
	lastUsedAd time.Time,
) *WebAuthnCredential {
	return &WebAuthnCredential{
		id:         id,
		credID:     credID,
		name:       name,
		publicKey:  publicKey,
		algorithm:  algorithm,
		transports: transports,
		createdAt:  createdAt,
		lastUsedAt: lastUsedAd,
	}
}

func (w *WebAuthnCredential) ID() values.WebAuthnCredentialID {
	return w.id
}

func (w *WebAuthnCredential) CredID() values.WebAuthnCredentialCredID {
	return w.credID
}

func (w *WebAuthnCredential) Name() values.WebAuthnCredentialName {
	return w.name
}

func (w *WebAuthnCredential) PublicKey() values.WebAuthnCredentialPublicKey {
	return w.publicKey
}

func (w *WebAuthnCredential) Algorithm() values.WebAuthnCredentialAlgorithm {
	return w.algorithm
}

func (w *WebAuthnCredential) Transports() []values.WebAuthnCredentialTransport {
	return w.transports
}

func (w *WebAuthnCredential) CreatedAt() time.Time {
	return w.createdAt
}

func (w *WebAuthnCredential) LastUsedAt() time.Time {
	return w.lastUsedAt
}

func (w *WebAuthnCredential) UpdateLastUsedAt() {
	w.lastUsedAt = time.Now()
}
