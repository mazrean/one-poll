package v1

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
	"github.com/mazrean/one-poll/service"
)

type WebAuthn struct {
	relyingPartyIDHash           values.WebAuthnRelyingPartyIDHash
	relyingParty                 *domain.WebAuthnRelyingParty
	db                           repository.DB
	webauthnCredentialRepository repository.WebAuthnCredential
	aaguid2NameMap               map[values.WebAuthnCredentialAAGUID]values.WebAuthnCredentialName
}

func NewWebAuthn(
	relyingParty *domain.WebAuthnRelyingParty,
	db repository.DB,
	webauthnCredentialRepository repository.WebAuthnCredential,
) (*WebAuthn, error) {
	aaguid2NameMap, err := loadAAGUID2NameMap()
	if err != nil {
		return nil, fmt.Errorf("failed to load aaguid2NameMap: %w", err)
	}

	return &WebAuthn{
		relyingPartyIDHash:           relyingParty.ID().Hash(),
		relyingParty:                 relyingParty,
		db:                           db,
		webauthnCredentialRepository: webauthnCredentialRepository,
		aaguid2NameMap:               aaguid2NameMap,
	}, nil
}

//go:generate curl -o webauthn_aaguid.json https://raw.githubusercontent.com/passkeydeveloper/passkey-authenticator-aaguids/main/aaguid.json
//go:embed webauthn_aaguid.json
var aaguidJSON []byte

func loadAAGUID2NameMap() (map[values.WebAuthnCredentialAAGUID]values.WebAuthnCredentialName, error) {
	var aaguidMap map[string]struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal(aaguidJSON, &aaguidMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal aaguid2NameMap: %w", err)
	}

	aaguid2NameMap := make(map[values.WebAuthnCredentialAAGUID]values.WebAuthnCredentialName, len(aaguidMap))
	for aaguidStr, v := range aaguidMap {
		aaguid, err := values.NewWebAuthnCredentialAAGUIDFromString(aaguidStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create aaguid: %w", err)
		}

		aaguid2NameMap[aaguid] = values.NewWebAuthnCredentialName(v.Name)
	}

	return aaguid2NameMap, nil
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
	credID values.WebAuthnCredentialCredID,
	aaguid values.WebAuthnCredentialAAGUID,
	publicKey values.WebAuthnCredentialPublicKey,
	algorithm values.WebAuthnCredentialAlgorithm,
) (*domain.WebAuthnCredential, error) {
	// ClientDataの検証
	if clientData.DataType() != values.WebAuthnClientDataTypeCreate {
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

	name, ok := wa.aaguid2NameMap[aaguid]
	if !ok {
		name = values.NewWebAuthnCredentialName("Unknown Authenticator")
	}

	now := time.Now()
	credential := domain.NewWebAuthnCredential(
		values.NewWebAuthnCredentialID(),
		credID,
		name,
		publicKey,
		algorithm,
		now,
		now,
	)

	err := wa.webauthnCredentialRepository.StoreCredential(ctx, user.GetID(), credential)
	if errors.Is(err, repository.ErrDuplicateRecord) {
		return nil, service.ErrWebAuthnDuplicate
	}
	if err != nil {
		return nil, fmt.Errorf("failed to store credential: %w", err)
	}

	return credential, nil
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
	if authData.RelyingPartyIDHash() != wa.relyingPartyIDHash {
		return nil, service.ErrWebAuthnInvalidRelyingParty
	}

	var user *domain.User
	err := wa.db.Transaction(ctx, nil, func(ctx context.Context) error {
		var (
			credential *domain.WebAuthnCredential
			err        error
		)
		credential, user, err = wa.webauthnCredentialRepository.GetCredentialWithUserByCredID(ctx, credID, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrWebAuthnInvalidCredential
		}
		if err != nil {
			return fmt.Errorf("failed to get credential: %w", err)
		}

		clientDataHash := clientData.Hash()

		switch credential.Algorithm() {
		case values.WebAuthnCredentialAlgorithmES256:
			err = wa.verifySignatureES256(
				append(authData.Raw(), clientDataHash[:]...),
				credential.PublicKey(),
				signature,
			)
		default:
			return fmt.Errorf("unsupported algorithm: %d", credential.Algorithm())
		}
		if err != nil {
			return fmt.Errorf("%w: %w", service.ErrWebAuthnInvalidSignature, err)
		}

		credential.UpdateLastUsedAt()
		err = wa.webauthnCredentialRepository.UpdateLastUsedAt(ctx, credential)
		// 同時刻に使用された場合は更新されないためErrNoRecordUpdatedは無視
		if err != nil && !errors.Is(err, repository.ErrNoRecordUpdated) {
			// 認証自体は成功しているため、ログだけ出して続行
			log.Printf("failed to update last used at: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return user, nil
}

func (wa *WebAuthn) verifySignatureES256(
	verificationData []byte,
	publicKey values.WebAuthnCredentialPublicKey,
	signature values.WebAuthnSignature,
) error {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKey)
	if x == nil {
		return errors.New("invalid public key")
	}

	ecdsaKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	hash := sha256.Sum256(verificationData)

	if !ecdsa.VerifyASN1(ecdsaKey, hash[:], signature) {
		return errors.New("invalid signature")
	}

	return nil
}

func (wa *WebAuthn) GetCredentials(ctx context.Context, user *domain.User) ([]*domain.WebAuthnCredential, error) {
	credentials, err := wa.webauthnCredentialRepository.GetCredentialsByUserID(ctx, user.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	return credentials, nil
}

func (wa *WebAuthn) DeleteCredential(ctx context.Context, user *domain.User, credID values.WebAuthnCredentialCredID) error {
	err := wa.webauthnCredentialRepository.DeleteCredential(ctx, user.GetID(), credID)
	if errors.Is(err, repository.ErrNoRecordDeleted) {
		return service.ErrWebAuthnNoCredential
	}
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	return nil
}
