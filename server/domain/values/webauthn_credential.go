package values

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type WebauthnCredentialID uuid.UUID

func NewWebauthnCredentialID() WebauthnCredentialID {
	return WebauthnCredentialID(uuid.New())
}

func NewWebauthnCredentialIDFromUUID(id uuid.UUID) WebauthnCredentialID {
	return WebauthnCredentialID(id)
}

type WebauthnCredentialCredID []byte

func NewWebauthnCredentialCredID(id []byte) WebauthnCredentialCredID {
	return WebauthnCredentialCredID(id)
}

var ErrWebauthnCredentialCredIDTooShort = errors.New("webauthn credential id should be at least 16 bytes")

func (id WebauthnCredentialCredID) Validate() error {
	if len(id) < 16 {
		return ErrWebauthnCredentialCredIDTooShort
	}

	return nil
}

type WebauthnCredentialName string

func NewWebauthnCredentialName(name string) WebauthnCredentialName {
	return WebauthnCredentialName(name)
}

var (
	ErrWebauthnCredentialNameEmpty   = errors.New("webauthn credential name is empty")
	ErrWebauthnCredentialNameTooLong = errors.New("webauthn credential name should be at most 64 characters")
)

func (name WebauthnCredentialName) Validate() error {
	if len(name) == 0 {
		return ErrWebauthnCredentialNameEmpty
	}

	length := utf8.RuneCountInString(string(name))
	if length > 64 {
		return ErrWebauthnCredentialNameTooLong
	}

	return nil
}

type WebauthnCredentialPublicKey []byte

func NewWebauthnCredentialPublicKey(key []byte) WebauthnCredentialPublicKey {
	return WebauthnCredentialPublicKey(key)
}

var ErrWebauthnCredentialPublicKeyEmpty = errors.New("webauthn credential public key is empty")

func (key WebauthnCredentialPublicKey) Validate() error {
	if len(key) == 0 {
		return ErrWebauthnCredentialPublicKeyEmpty
	}

	return nil
}

type WebauthnCredentialAlgorithm uint

const (
	WebauthnCredentialAlgorithmES256 WebauthnCredentialAlgorithm = iota + 1
)

type WebauthnCredentialTransport uint8

const (
	WebauthnCredentialTransportUSB WebauthnCredentialTransport = iota + 1
	WebauthnCredentialTransportNFC
	WebauthnCredentialTransportBLE
	WebauthnCredentialTransportSmartCard
	WebauthnCredentialTransportHybrid
	WebAuthnCredentialTransportInternal
)

func (t WebauthnCredentialTransport) String() string {
	switch t {
	case WebauthnCredentialTransportUSB:
		return "usb"
	case WebauthnCredentialTransportNFC:
		return "nfc"
	case WebauthnCredentialTransportBLE:
		return "ble"
	case WebauthnCredentialTransportSmartCard:
		return "smart_card"
	case WebauthnCredentialTransportHybrid:
		return "hybrid"
	case WebAuthnCredentialTransportInternal:
		return "internal"
	default:
		return ""
	}
}
