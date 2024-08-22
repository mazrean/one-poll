package values

import (
	"crypto/sha256"
	"errors"
	"regexp"
)

type WebAuthnRelyingPartyID string

func NewWebAuthnRelyingPartyID(id string) WebAuthnRelyingPartyID {
	return WebAuthnRelyingPartyID(id)
}

var (
	ErrWebAuthnRelyingPartyIDInvalidDomain = errors.New("invalid domain")
)

var domainRegexp = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)

func (id WebAuthnRelyingPartyID) Validate() error {
	if !domainRegexp.MatchString(string(id)) && string(id) != "localhost" {
		return ErrWebAuthnRelyingPartyIDInvalidDomain
	}

	return nil
}

func (id WebAuthnRelyingPartyID) Hash() WebAuthnRelyingPartyIDHash {
	return WebAuthnRelyingPartyIDHash(sha256.Sum256([]byte(id)))
}

type WebAuthnRelyingPartyDisplayName string

func NewWebAuthnRelyingPartyDisplayName(name string) WebAuthnRelyingPartyDisplayName {
	return WebAuthnRelyingPartyDisplayName(name)
}
