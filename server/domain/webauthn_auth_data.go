package domain

import "github.com/mazrean/one-poll/domain/values"

type WebAuthnAuthData struct {
	relyingPartyIDHash values.WebAuthnRelyingPartyIDHash
	raw                []byte
}

func NewWebAuthnAuthData(
	rpIDHash values.WebAuthnRelyingPartyIDHash,
	raw []byte,
) *WebAuthnAuthData {
	return &WebAuthnAuthData{relyingPartyIDHash: rpIDHash}
}

func (w *WebAuthnAuthData) RelyingPartyIDHash() values.WebAuthnRelyingPartyIDHash {
	return w.relyingPartyIDHash
}

func (w *WebAuthnAuthData) Raw() []byte {
	return w.raw
}
