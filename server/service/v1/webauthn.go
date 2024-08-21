package v1

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
	"github.com/mazrean/one-poll/service"
)

type WebAuthn struct {
	relyingPartyIDHash           values.WebAuthnRelyingPartyIDHash
	relyingParty                 *domain.WebAuthnRelyingParty
	db                           repository.DB
	webauthnCredentialRepository repository.WebauthnCredential
}

func NewWebAuthn(
	relyingParty *domain.WebAuthnRelyingParty,
	db repository.DB,
	webauthnCredentialRepository repository.WebauthnCredential,
) *WebAuthn {
	return &WebAuthn{
		relyingPartyIDHash:           relyingParty.ID().Hash(),
		relyingParty:                 relyingParty,
		db:                           db,
		webauthnCredentialRepository: webauthnCredentialRepository,
	}
}

func (wa *WebAuthn) BeginRegistration(ctx context.Context, _ *domain.User) (*domain.WebAuthnRelyingParty, values.WebAuthnChallenge, error) {
	challenge, err := values.NewWebAuthnChallenge()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	return wa.relyingParty, challenge, nil
}

func (wa *WebAuthn) FinishRegistration(
	ctx context.Context,
	user *domain.User,
	sessionChallenge values.WebAuthnChallenge,
	relyingPartyHash values.WebAuthnRelyingPartyIDHash,
	clientData *domain.WebAuthnClientData,
	credential *domain.WebAuthnCredential,
) error {
	// ClientDataの検証
	if clientData.DataType() != values.WebAuthnClientDataTypeCreate {
		return service.ErrWebAuthnInvalidDataType
	}

	if sessionChallenge.ConstantTimeEqual(clientData.Challenge()) {
		return service.ErrWebAuthnInvalidChallenge
	}

	if clientData.Origin() != wa.relyingParty.Origin() {
		return service.ErrWebAuthnInvalidOrigin
	}

	// RelyingPartyIDHashの検証
	if relyingPartyHash != wa.relyingPartyIDHash {
		return service.ErrWebAuthnInvalidRelyingParty
	}

	err := wa.webauthnCredentialRepository.StoreCredential(ctx, user.GetID(), credential)
	if errors.Is(err, repository.ErrDuplicateRecord) {
		return service.ErrWebAuthnDuplicate
	}
	if err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}

	return nil
}

func (wa *WebAuthn) BeginLogin(ctx context.Context) (*domain.WebAuthnRelyingParty, values.WebAuthnChallenge, error) {
	challenge, err := values.NewWebAuthnChallenge()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	return wa.relyingParty, challenge, nil
}

func (wa *WebAuthn) FinishLogin(
	ctx context.Context,
	sessionChallenge values.WebAuthnChallenge,
	relyingPartyHash values.WebAuthnRelyingPartyIDHash,
	clientData *domain.WebAuthnClientData,
	authData *domain.WebAuthnAuthData,
	credID values.WebAuthnCredentialCredID,
	signature values.WebAuthnSignature,
) (*domain.User, error) {
	// ClientDataの検証
	if clientData.DataType() != values.WebAuthnClientDataTypeGet {
		return nil, service.ErrWebAuthnInvalidDataType
	}

	if sessionChallenge.ConstantTimeEqual(clientData.Challenge()) {
		return nil, service.ErrWebAuthnInvalidChallenge
	}

	if clientData.Origin() != wa.relyingParty.Origin() {
		return nil, service.ErrWebAuthnInvalidOrigin
	}

	// RelyingPartyIDHashの検証
	if relyingPartyHash != wa.relyingPartyIDHash {
		return nil, service.ErrWebAuthnInvalidRelyingParty
	}

	credential, user, err := wa.webauthnCredentialRepository.GetCredentialWithUserByCredID(ctx, credID)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrWebAuthnInvalidCredential
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	verificationData := bytes.NewBuffer(nil)
	verificationData.Write(authData.Raw())
	clientDataHash := clientData.Hash()
	verificationData.Write(clientDataHash[:])

	switch credential.Algorithm() {
	case values.WebAuthnCredentialAlgorithmES256:
		err = wa.verifySignatureES256(verificationData.Bytes(), credential.PublicKey(), signature)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %d", credential.Algorithm())
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", service.ErrWebAuthnInvalidSignature, err)
	}

	return user, nil
}

func (wa *WebAuthn) verifySignatureES256(
	verificationData []byte,
	publicKey values.WebAuthnCredentialPublicKey,
	sig values.WebAuthnSignature,
) error {
	key, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	ecdsaKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("not an ECDSA public key")
	}

	signature := &struct {
		R, S *big.Int
	}{}
	_, err = asn1.Unmarshal(sig, signature)
	if err != nil {
		return nil
	}

	h := sha256.New()
	h.Write(verificationData)

	if !ecdsa.Verify(ecdsaKey, h.Sum(nil), signature.R, signature.S) {
		return errors.New("invalid signature")
	}

	return nil
}
