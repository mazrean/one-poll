package domain

import "github.com/mazrean/one-poll/domain/values"

type WebAuthnAuthData struct {
	relyingPartyIDHash values.WebAuthnRelyingPartyIDHash
}

func NewWebAuthnAuthData(rpIDHash values.WebAuthnRelyingPartyIDHash) *WebAuthnAuthData {
	return &WebAuthnAuthData{relyingPartyIDHash: rpIDHash}
}

func (w *WebAuthnAuthData) RelyingPartyIDHash() values.WebAuthnRelyingPartyIDHash {
	return w.relyingPartyIDHash
}
