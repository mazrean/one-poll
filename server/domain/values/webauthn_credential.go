package values

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type WebAuthnCredentialID uuid.UUID

func NewWebauthnCredentialID() WebAuthnCredentialID {
	return WebAuthnCredentialID(uuid.New())
}

func NewWebauthnCredentialIDFromUUID(id uuid.UUID) WebAuthnCredentialID {
	return WebAuthnCredentialID(id)
}

type WebAuthnCredentialCredID []byte

func NewWebauthnCredentialCredID(id []byte) WebAuthnCredentialCredID {
	return WebAuthnCredentialCredID(id)
}

var ErrWebauthnCredentialCredIDTooShort = errors.New("webauthn credential id should be at least 16 bytes")

func (id WebAuthnCredentialCredID) Validate() error {
	if len(id) < 16 {
		return ErrWebauthnCredentialCredIDTooShort
	}

	return nil
}

type WebAuthnCredentialName string

func NewWebauthnCredentialName(name string) WebAuthnCredentialName {
	return WebAuthnCredentialName(name)
}

var (
	ErrWebauthnCredentialNameEmpty   = errors.New("webauthn credential name is empty")
	ErrWebauthnCredentialNameTooLong = errors.New("webauthn credential name should be at most 64 characters")
)

func (name WebAuthnCredentialName) Validate() error {
	if len(name) == 0 {
		return ErrWebauthnCredentialNameEmpty
	}

	length := utf8.RuneCountInString(string(name))
	if length > 64 {
		return ErrWebauthnCredentialNameTooLong
	}

	return nil
}

type WebAuthnCredentialPublicKey []byte

func NewWebauthnCredentialPublicKey(key []byte) WebAuthnCredentialPublicKey {
	return WebAuthnCredentialPublicKey(key)
}

var ErrWebauthnCredentialPublicKeyEmpty = errors.New("webauthn credential public key is empty")

func (key WebAuthnCredentialPublicKey) Validate() error {
	if len(key) == 0 {
		return ErrWebauthnCredentialPublicKeyEmpty
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
	WebauthnCredentialTransportUSB WebAuthnCredentialTransport = iota + 1
	WebauthnCredentialTransportNFC
	WebauthnCredentialTransportBLE
	WebauthnCredentialTransportSmartCard
	WebauthnCredentialTransportHybrid
	WebAuthnCredentialTransportInternal
)

func (t WebAuthnCredentialTransport) String() string {
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
