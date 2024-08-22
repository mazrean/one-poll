package domain

import "github.com/mazrean/one-poll/domain/values"

type WebAuthnRelyingParty struct {
	id          values.WebAuthnRelyingPartyID
	origin      values.WebAuthnOrigin
	displayName values.WebAuthnRelyingPartyDisplayName
}

func NewWebAuthnRelyingParty(
	id values.WebAuthnRelyingPartyID,
	origin values.WebAuthnOrigin,
	displayName values.WebAuthnRelyingPartyDisplayName,
) *WebAuthnRelyingParty {
	return &WebAuthnRelyingParty{
		id:          id,
		origin:      origin,
		displayName: displayName,
	}
}

func (w *WebAuthnRelyingParty) ID() values.WebAuthnRelyingPartyID {
	return w.id
}

func (w *WebAuthnRelyingParty) Origin() values.WebAuthnOrigin {
	return w.origin
}

func (w *WebAuthnRelyingParty) DisplayName() values.WebAuthnRelyingPartyDisplayName {
	return w.displayName
}
