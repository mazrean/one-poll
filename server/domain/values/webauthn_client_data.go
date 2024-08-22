package values

import (
	"crypto/sha256"
)

type WebAuthnClientDataType uint8

const (
	WebAuthnClientDataTypeCreate WebAuthnClientDataType = iota + 1
	WebAuthnClientDataTypeGet
)

func (t WebAuthnClientDataType) String() string {
	switch t {
	case WebAuthnClientDataTypeCreate:
		return "create"
	case WebAuthnClientDataTypeGet:
		return "get"
	default:
		return "unknown"
	}
}

type WebAuthnClientDataHash [32]byte

func NewWebAuthnClientDataHash(hash [32]byte) WebAuthnClientDataHash {
	return WebAuthnClientDataHash(hash)
}

func NewWebAuthnClientDataHashFromRaw(raw []byte) WebAuthnClientDataHash {
	return WebAuthnClientDataHash(sha256.Sum256(raw))
}
