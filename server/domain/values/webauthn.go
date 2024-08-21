package values

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"github.com/mazrean/one-poll/pkg/random"
)

type WebAuthnChallenge []byte

func NewWebAuthnChallenge() (WebAuthnChallenge, error) {
	challenge, err := random.Secure(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate challenge: %v", err)
	}

	return WebAuthnChallenge(challenge), nil
}

var ErrWebAuthnSessionChallengeInvalid = errors.New("challenge is invalid")

func NewWebAuthnChallengeFromBase64URL(challenge string) (WebAuthnChallenge, error) {
	b, err := base64.URLEncoding.DecodeString(challenge)
	if err != nil {
		return nil, ErrWebAuthnSessionChallengeInvalid
	}

	return WebAuthnChallenge(b), nil
}

var ErrWebAuthnSessionChallengeInvalidLength = errors.New("challenge length is invalid")

func (c WebAuthnChallenge) Validate() error {
	if len(c) != 32 {
		return ErrWebAuthnSessionChallengeInvalidLength
	}

	return nil
}

func (c WebAuthnChallenge) ConstantTimeEqual(other WebAuthnChallenge) bool {
	return subtle.ConstantTimeCompare(c, other) == 1
}

type WebAuthnOrigin string

func NewWebAuthnOrigin(origin string) WebAuthnOrigin {
	return WebAuthnOrigin(origin)
}

var (
	ErrWebAuthnOriginInvalidURL       = errors.New("invalid url")
	ErrWebAuthnOriginNotHTTPS         = errors.New("scheme is not https")
	ErrWebAuthnOriginHostEmpty        = errors.New("host is empty")
	ErrWebAuthnOriginPathNotEmpty     = errors.New("path is not empty")
	ErrWebAuthnOriginQueryNotEmpty    = errors.New("query is not empty")
	ErrWebAuthnOriginFragmentNotEmpty = errors.New("fragment is not empty")
	ErrWebAuthnOriginUserNotNil       = errors.New("user is not nil")
)

func (origin WebAuthnOrigin) Validate() error {
	originURL, err := url.Parse(string(origin))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWebAuthnOriginInvalidURL, err)
	}

	if originURL.Scheme != "https" {
		return fmt.Errorf("%w: %s", ErrWebAuthnOriginNotHTTPS, originURL.Scheme)
	}

	if originURL.Host == "" {
		return ErrWebAuthnOriginHostEmpty
	}

	if originURL.RawPath != "" {
		return fmt.Errorf("%w: %s", ErrWebAuthnOriginPathNotEmpty, originURL.RawPath)
	}

	if originURL.RawQuery != "" {
		return fmt.Errorf("%w: %s", ErrWebAuthnOriginQueryNotEmpty, originURL.RawQuery)
	}

	if originURL.Fragment != "" {
		return fmt.Errorf("%w: %s", ErrWebAuthnOriginFragmentNotEmpty, originURL.Fragment)
	}

	if originURL.User != nil {
		return ErrWebAuthnOriginUserNotNil
	}

	return nil
}

type WebAuthnRelyingPartyIDHash [32]byte

func NewWebAuthnRelyingPartyIDHash(hash [32]byte) WebAuthnRelyingPartyIDHash {
	return WebAuthnRelyingPartyIDHash(hash)
}

type WebAuthnSignature []byte

func NewWebAuthnSignature(signature []byte) WebAuthnSignature {
	return WebAuthnSignature(signature)
}
