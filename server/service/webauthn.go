package service

import (
	"context"
	"errors"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
)

var (
	ErrWebAuthnInvalidDataType     = errors.New("invalid data type")
	ErrWebAuthnInvalidChallenge    = errors.New("invalid challenge")
	ErrWebAuthnInvalidOrigin       = errors.New("invalid origin")
	ErrWebAuthnInvalidRelyingParty = errors.New("invalid relying party")
	ErrWebAuthnDuplicate           = errors.New("duplicate credential")
	ErrWebAuthnInvalidSignature    = errors.New("invalid signature")
	ErrWebAuthnInvalidCredential   = errors.New("invalid credential")
)

type WebAuthn interface {
	// BeginRegistration クレデンシャル登録の開始
	BeginRegistration(ctx context.Context, user *domain.User) (*domain.WebAuthnRelyingParty, values.WebAuthnChallenge, error)
	/*
		FinishRegistration clientDataを検証し、クレデンシャルを保存する
		エラー: ErrWebAuthnInvalidDataType, ErrWebAuthnInvalidChallenge, ErrWebAuthnInvalidOrigin, ErrWebAuthnInvalidRelyingParty, ErrWebAuthnDuplicate
	*/
	FinishRegistration(
		ctx context.Context,
		user *domain.User,
		sessionChallenge values.WebAuthnChallenge,
		relyingPartyHash values.WebAuthnRelyingPartyIDHash,
		clientData *domain.WebAuthnClientData,
		credID values.WebAuthnCredentialCredID,
		aaguid values.WebAuthnCredentialAAGUID,
		publicKey values.WebAuthnCredentialPublicKey,
		algorithm values.WebAuthnCredentialAlgorithm,
		transports []values.WebAuthnCredentialTransport,
	) (*domain.WebAuthnCredential, error)
	// BeginLogin ログインの開始
	BeginLogin(ctx context.Context) (*domain.WebAuthnRelyingParty, values.WebAuthnChallenge, error)
	/*
		FinishLogin clientData、signatureを検証し、ログインを完了する
		エラー: ErrWebAuthnInvalidDataType, ErrWebAuthnInvalidChallenge, ErrWebAuthnInvalidOrigin, ErrWebAuthnInvalidRelyingParty, ErrWebAuthnInvalidCredential, ErrWebAuthnInvalidSignature
	*/
	FinishLogin(
		ctx context.Context,
		sessionChallenge values.WebAuthnChallenge,
		relyingPartyHash values.WebAuthnRelyingPartyIDHash,
		clientData *domain.WebAuthnClientData,
		authData *domain.WebAuthnAuthData,
		credID values.WebAuthnCredentialCredID,
		signature values.WebAuthnSignature,
	) (*domain.User, error)
	DeleteCredential(ctx context.Context, user *domain.User, credID values.WebAuthnCredentialCredID) error
}
