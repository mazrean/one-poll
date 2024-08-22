package v1

import (
	"crypto/elliptic"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/handler/v1/openapi"
	"github.com/mazrean/one-poll/service"
	"github.com/ugorji/go/codec"
)

type WebAuthn struct {
	*Session
	webAuthnService  service.WebAuthn
	pubKeyCredParams []openapi.WebAuthnPubKeyCredParam
}

func NewWebAuthn(session *Session, webAuthnService service.WebAuthn) *WebAuthn {
	var pubKeyCredParams []openapi.WebAuthnPubKeyCredParam
	for _, algorism := range values.WebAuthnCredentialAlgorithms {
		switch algorism {
		case values.WebAuthnCredentialAlgorithmES256:
			pubKeyCredParams = append(pubKeyCredParams, openapi.WebAuthnPubKeyCredParam{
				Type: "public-key",
				Alg:  -7,
			})
		default:
			log.Printf("unsupported algorithm: %v", algorism)
		}
	}

	return &WebAuthn{
		Session:          session,
		webAuthnService:  webAuthnService,
		pubKeyCredParams: pubKeyCredParams,
	}
}

func (w *WebAuthn) GetWebauthnCredentials(c echo.Context) error {
	session, err := w.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := w.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil && !errors.Is(err, ErrNoValue) {
		log.Printf("failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	credentials, err := w.webAuthnService.GetCredentials(c.Request().Context(), user)
	if err != nil {
		log.Printf("failed to get credentials: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get credentials")
	}

	apiCredentials := make([]openapi.WebAuthnCredential, 0, len(credentials))
	for _, credential := range credentials {
		apiCredentials = append(apiCredentials, openapi.WebAuthnCredential{
			Id:         uuid.UUID(credential.ID()),
			Name:       string(credential.Name()),
			CreatedAt:  credential.CreatedAt(),
			LastUsedAt: credential.LastUsedAt(),
		})
	}

	return c.JSON(http.StatusOK, apiCredentials)
}

// webauthnの公開鍵登録開始
// (POST /webauthn/resister/start)
func (w *WebAuthn) PostWebauthnResisterStart(c echo.Context) error {
	session, err := w.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := w.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil && !errors.Is(err, ErrNoValue) {
		log.Printf("failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	relyingParty, challenge, err := w.webAuthnService.BeginRegistration(c.Request().Context(), user)
	if err != nil {
		log.Printf("failed to start registration: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start registration")
	}

	w.Session.setWebAuthnResisterChallenge(session, challenge)

	err = w.Session.save(c, session)
	if err != nil {
		log.Printf("failed to save session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save session")
	}

	relyingPartyID := string(relyingParty.ID())
	strChallenge := base64.RawURLEncoding.EncodeToString(challenge)

	requireResidentKey := true
	residentKey := openapi.Required

	return c.JSON(http.StatusOK, openapi.WebAuthnPublicKeyCredentialCreationOptions{
		Rp: openapi.WebAuthnRelyingParty{
			Id:   &relyingPartyID,
			Name: "One Poll",
		},
		User: openapi.WebAuthnUser{
			Id:          uuid.UUID(user.GetID()),
			Name:        string(user.GetName()),
			DisplayName: string(user.GetName()),
		},
		Timeout:          60000,
		Challenge:        strChallenge,
		PubKeyCredParams: w.pubKeyCredParams,
		Attestation:      openapi.Direct,
		AuthenticatorSelection: openapi.WebAuthnAuthenticatorSelectionCriteria{
			RequireResidentKey: &requireResidentKey,
			ResidentKey:        &residentKey,
		},
	})
}

// webauthnの公開鍵登録終了
// (POST /webauthn/resister/finish)
func (w *WebAuthn) PostWebauthnResisterFinish(c echo.Context) error {
	session, err := w.Session.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	}

	user, err := w.Session.getUser(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusUnauthorized, "login required")
	}
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	challenge, err := w.Session.getWebAuthnRegisterChallenge(session)
	if errors.Is(err, ErrNoValue) {
		return echo.NewHTTPError(http.StatusBadRequest, "register challenge not found")
	}
	if err != nil {
		log.Printf("failed to get register challenge: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get register challenge")
	}

	var req openapi.WebAuthnPublicKeyCredentialCreation
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	clientData, err := parseClientData(req.Response.ClientDataJSON)
	if err != nil {
		log.Printf("failed to parse client data: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid client data")
	}

	attestationObject, err := parseAttestationObject(req.Response.AttestationObject)
	if err != nil {
		log.Printf("failed to parse attestation object: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid attestation object")
	}

	authData, err := parseAuthData(attestationObject.AuthData)
	if err != nil {
		log.Printf("failed to parse auth data: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid auth data")
	}

	publicKey, algorism, err := parsePublicKey(authData.PublicKey)
	if err != nil {
		log.Printf("failed to parse public key: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid public key")
	}

	credential, err := w.webAuthnService.FinishRegistration(
		c.Request().Context(),
		user,
		challenge,
		authData.RPIDHash,
		clientData,
		authData.CredID,
		authData.AAGUID,
		publicKey,
		algorism,
	)
	if err != nil {
		log.Printf("failed to finish registration: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to finish registration")
	}

	return c.JSON(http.StatusOK, openapi.WebAuthnCredential{
		Id:         uuid.UUID(credential.ID()),
		Name:       string(credential.Name()),
		CreatedAt:  credential.CreatedAt(),
		LastUsedAt: credential.LastUsedAt(),
	})
}

type attestationObject struct {
	AuthData []byte `codec:"authData" cbor:"authData"`
}

func parseAttestationObject(req string) (*attestationObject, error) {
	attestationObjectBytes, err := base64.RawURLEncoding.DecodeString(req)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestation object: %w", err)
	}

	var obj attestationObject
	err = codec.NewDecoderBytes(attestationObjectBytes, new(codec.CborHandle)).Decode(&obj)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestation object: %w", err)
	}

	return &obj, nil
}

type authenticatorData struct {
	RPIDHash  values.WebAuthnRelyingPartyIDHash
	AAGUID    values.WebAuthnCredentialAAGUID
	CredID    values.WebAuthnCredentialCredID
	PublicKey []byte
}

func parseAuthData(authData []byte) (*authenticatorData, error) {
	rpIDHash := [32]byte(authData[0:32])
	aaguid := [16]byte(authData[37:53])
	credentialIDLength := binary.BigEndian.Uint16(authData[53:55])
	credentialID := authData[55 : 55+credentialIDLength]
	credentialPublicKey := authData[55+credentialIDLength:]

	return &authenticatorData{
		RPIDHash:  values.NewWebAuthnRelyingPartyIDHash(rpIDHash),
		AAGUID:    values.NewWebAuthnCredentialAAGUID(aaguid),
		CredID:    values.NewWebAuthnCredentialCredID(credentialID),
		PublicKey: credentialPublicKey,
	}, nil
}

func parsePublicKey(publicKey []byte) (values.WebAuthnCredentialPublicKey, values.WebAuthnCredentialAlgorithm, error) {
	var cborMap map[int]any
	err := codec.NewDecoderBytes(publicKey, new(codec.CborHandle)).Decode(&cborMap)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode map: %w", err)
	}

	alg, ok := cborMap[3].(int64)
	if !ok {
		return nil, 0, fmt.Errorf("failed to get alg: %t", cborMap[3])
	}

	var (
		pubKey   values.WebAuthnCredentialPublicKey
		algorism values.WebAuthnCredentialAlgorithm
	)
	switch alg {
	case -7:
		xBinary, ok := cborMap[-2].([]byte)
		if !ok {
			return nil, 0, fmt.Errorf("failed to get x: %t", cborMap[-2])
		}

		yBinary, ok := cborMap[-3].([]byte)
		if !ok {
			return nil, 0, fmt.Errorf("failed to get y: %t", cborMap[-3])
		}

		if len(xBinary) != 32 || len(yBinary) != 32 {
			return nil, 0, fmt.Errorf("invalid public key length(X: %d, Y: %d)", len(xBinary), len(yBinary))
		}

		x := big.NewInt(0)
		x.SetBytes(xBinary)

		y := big.NewInt(0)
		y.SetBytes(yBinary)

		publicKeyBuffer := elliptic.MarshalCompressed(
			elliptic.P256(),
			x,
			y,
		)

		pubKey = values.NewWebAuthnCredentialPublicKey(publicKeyBuffer)
		algorism = values.WebAuthnCredentialAlgorithmES256
	default:
		return nil, 0, fmt.Errorf("unsupported algorithm: %d", alg)
	}

	return pubKey, algorism, nil
}

func parseClientData(clientDataJSON string) (*domain.WebAuthnClientData, error) {
	clientDataBytes, err := base64.RawURLEncoding.DecodeString(clientDataJSON)
	if err != nil {
		return nil, err
	}

	var reqClientData struct {
		Type      string `json:"type"`
		Challenge string `json:"challenge"`
		Origin    string `json:"origin"`
	}
	if err := json.Unmarshal(clientDataBytes, &reqClientData); err != nil {
		return nil, err
	}

	var dataType values.WebAuthnClientDataType
	switch reqClientData.Type {
	case "webauthn.create":
		dataType = values.WebAuthnClientDataTypeCreate
	case "webauthn.get":
		dataType = values.WebAuthnClientDataTypeGet
	default:
		return nil, errors.New("unsupported client data type")
	}

	challenge, err := values.NewWebAuthnChallengeFromBase64URL(reqClientData.Challenge)
	if err != nil {
		return nil, fmt.Errorf("invalid challenge: %w", err)
	}

	return domain.NewWebAuthnClientData(
		dataType,
		challenge,
		values.NewWebAuthnOrigin(reqClientData.Origin),
		values.NewWebAuthnClientDataHashFromRaw(clientDataBytes),
	), nil
}
