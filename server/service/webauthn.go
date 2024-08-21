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
		エラー: ErrWebAuthnInvalidChallenge, ErrWebAuthnInvalidOrigin, ErrWebAuthnDuplicate
	*/
	FinishRegistration(
		ctx context.Context,
		user *domain.User,
		sessionChallenge values.WebAuthnChallenge,
		relyingPartyHash values.WebAuthnRelyingPartyIDHash,
		clientData *domain.WebAuthnClientData,
		credential *domain.WebAuthnCredential,
	) error
	// BeginLogin ログインの開始
	BeginLogin(ctx context.Context) (*domain.WebAuthnRelyingParty, values.WebAuthnChallenge, error)
	/*
		FinishLogin clientData、signatureを検証し、ログインを完了する
		エラー: ErrWebAuthnInvalidChallenge, ErrWebAuthnInvalidOrigin, ErrWebAuthnInvalidSignature
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
}
