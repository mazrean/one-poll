package values

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type WebAuthnCredentialID uuid.UUID

func NewWebAuthnCredentialID() WebAuthnCredentialID {
	return WebAuthnCredentialID(uuid.New())
}

func NewWebAuthnCredentialIDFromUUID(id uuid.UUID) WebAuthnCredentialID {
	return WebAuthnCredentialID(id)
}

type WebAuthnCredentialCredID []byte

func NewWebAuthnCredentialCredID(id []byte) WebAuthnCredentialCredID {
	return WebAuthnCredentialCredID(id)
}

var ErrWebAuthnCredentialCredIDTooShort = errors.New("webauthn credential id should be at least 16 bytes")

func (id WebAuthnCredentialCredID) Validate() error {
	if len(id) < 16 {
		return ErrWebAuthnCredentialCredIDTooShort
	}

	return nil
}

type WebAuthnCredentialAAGUID uuid.UUID

func NewWebAuthnCredentialAAGUID(guid uuid.UUID) WebAuthnCredentialAAGUID {
	return WebAuthnCredentialAAGUID(guid)
}

func NewWebAuthnCredentialAAGUIDFromString(guid string) (WebAuthnCredentialAAGUID, error) {
	id, err := uuid.Parse(guid)
	if err != nil {
		return WebAuthnCredentialAAGUID{}, err
	}

	return WebAuthnCredentialAAGUID(id), nil
}

func (guid WebAuthnCredentialAAGUID) String() string {
	return uuid.UUID(guid).String()
}

type WebAuthnCredentialName string

func NewWebAuthnCredentialName(name string) WebAuthnCredentialName {
	return WebAuthnCredentialName(name)
}

var (
	ErrWebAuthnCredentialNameEmpty   = errors.New("webauthn credential name is empty")
	ErrWebAuthnCredentialNameTooLong = errors.New("webauthn credential name should be at most 64 characters")
)

func (name WebAuthnCredentialName) Validate() error {
	if len(name) == 0 {
		return ErrWebAuthnCredentialNameEmpty
	}

	length := utf8.RuneCountInString(string(name))
	if length > 64 {
		return ErrWebAuthnCredentialNameTooLong
	}

	return nil
}

type WebAuthnCredentialPublicKey []byte

func NewWebAuthnCredentialPublicKey(key []byte) WebAuthnCredentialPublicKey {
	return WebAuthnCredentialPublicKey(key)
}

var ErrWebAuthnCredentialPublicKeyEmpty = errors.New("webauthn credential public key is empty")

func (key WebAuthnCredentialPublicKey) Validate() error {
	if len(key) == 0 {
		return ErrWebAuthnCredentialPublicKeyEmpty
	}

	return nil
}

type WebAuthnCredentialAlgorithm uint

const (
	WebAuthnCredentialAlgorithmES256 WebAuthnCredentialAlgorithm = iota + 1
)

var WebAuthnCredentialAlgorithms = []WebAuthnCredentialAlgorithm{
	WebAuthnCredentialAlgorithmES256,
}

type WebAuthnCredentialTransport uint8

const (
	WebAuthnCredentialTransportUSB WebAuthnCredentialTransport = iota + 1
	WebAuthnCredentialTransportNFC
	WebAuthnCredentialTransportBLE
	WebAuthnCredentialTransportSmartCard
	WebAuthnCredentialTransportHybrid
	WebAuthnCredentialTransportInternal
)

func (t WebAuthnCredentialTransport) String() string {
	switch t {
	case WebAuthnCredentialTransportUSB:
		return "usb"
	case WebAuthnCredentialTransportNFC:
		return "nfc"
	case WebAuthnCredentialTransportBLE:
		return "ble"
	case WebAuthnCredentialTransportSmartCard:
		return "smart_card"
	case WebAuthnCredentialTransportHybrid:
		return "hybrid"
	case WebAuthnCredentialTransportInternal:
		return "internal"
	default:
		return ""
	}
}
