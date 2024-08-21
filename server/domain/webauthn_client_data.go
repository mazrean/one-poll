package domain

import "github.com/mazrean/one-poll/domain/values"

type WebAuthnClientData struct {
	dataType  values.WebAuthnClientDataType
	challenge values.WebAuthnChallenge
	origin    values.WebAuthnOrigin
	hash      values.WebAuthnClientDataHash
}

func NewWebAuthnClientData(
	dataType values.WebAuthnClientDataType,
	challenge values.WebAuthnChallenge,
	origin values.WebAuthnOrigin,
	hash values.WebAuthnClientDataHash,
) *WebAuthnClientData {
	return &WebAuthnClientData{
		dataType:  dataType,
		challenge: challenge,
		origin:    origin,
		hash:      hash,
	}
}

func (w *WebAuthnClientData) DataType() values.WebAuthnClientDataType {
	return w.dataType
}

func (w *WebAuthnClientData) Challenge() values.WebAuthnChallenge {
	return w.challenge
}

func (w *WebAuthnClientData) Origin() values.WebAuthnOrigin {
	return w.origin
}

func (w *WebAuthnClientData) Hash() values.WebAuthnClientDataHash {
	return w.hash
}
