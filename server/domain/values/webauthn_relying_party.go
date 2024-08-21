package values

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
)

type WebAuthnRelyingPartyID string

func NewWebAuthnRelyingPartyID(id string) WebAuthnRelyingPartyID {
	return WebAuthnRelyingPartyID(id)
}

var (
	ErrWebAuthnRelyingPartyIDInvalidURL       = errors.New("invalid url")
	ErrWebAuthnRelyingPartyIDSchemeExists     = errors.New("scheme exists")
	ErrWebAuthnRelyingPartyIDHostEmpty        = errors.New("host is empty")
	ErrWebAuthnRelyingPartyIDPathNotEmpty     = errors.New("path is not empty")
	ErrWebAuthnRelyingPartyIDQueryNotEmpty    = errors.New("query is not empty")
	ErrWebAuthnRelyingPartyIDFragmentNotEmpty = errors.New("fragment is not empty")
	ErrWebAuthnRelyingPartyIDUserNotNil       = errors.New("user is not nil")
)

func (id WebAuthnRelyingPartyID) Validate() error {
	idURL, err := url.Parse(string(id))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWebAuthnRelyingPartyIDInvalidURL, err)
	}

	if idURL.Scheme != "" {
		return ErrWebAuthnRelyingPartyIDSchemeExists
	}

	if idURL.Host == "" {
		return ErrWebAuthnRelyingPartyIDHostEmpty
	}

	if idURL.RawPath != "" {
		return fmt.Errorf("%w: %s", ErrWebAuthnRelyingPartyIDPathNotEmpty, idURL.RawPath)
	}

	if idURL.RawQuery != "" {
		return fmt.Errorf("%w: %s", ErrWebAuthnRelyingPartyIDQueryNotEmpty, idURL.RawQuery)
	}

	if idURL.Fragment != "" {
		return fmt.Errorf("%w: %s", ErrWebAuthnRelyingPartyIDFragmentNotEmpty, idURL.Fragment)
	}

	if idURL.User != nil {
		return ErrWebAuthnRelyingPartyIDUserNotNil
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
