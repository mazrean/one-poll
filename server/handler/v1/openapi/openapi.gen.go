// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for PollStatus.
const (
	Limited  PollStatus = "limited"
	Opened   PollStatus = "opened"
	Outdated PollStatus = "outdated"
)

// Defines values for PollType.
const (
	Radio PollType = "radio"
)

// Defines values for UserStatusAccessMode.
const (
	CanAccessDetails UserStatusAccessMode = "can_access_details"
	CanAnswer        UserStatusAccessMode = "can_answer"
	OnlyBrowsable    UserStatusAccessMode = "only_browsable"
)

// Defines values for WebAuthnAuthenticatorAttachment.
const (
	CrossPlatform WebAuthnAuthenticatorAttachment = "cross-platform"
	Platform      WebAuthnAuthenticatorAttachment = "platform"
)

// Defines values for WebAuthnAuthenticatorAttestationType.
const (
	Direct   WebAuthnAuthenticatorAttestationType = "direct"
	Indirect WebAuthnAuthenticatorAttestationType = "indirect"
	None     WebAuthnAuthenticatorAttestationType = "none"
)

// Defines values for WebAuthnAuthenticatorResidentKeyRequirement.
const (
	Discouraged WebAuthnAuthenticatorResidentKeyRequirement = "discouraged"
	Preferred   WebAuthnAuthenticatorResidentKeyRequirement = "preferred"
	Required    WebAuthnAuthenticatorResidentKeyRequirement = "required"
)

// Defines values for WebAuthnCredentialAlgorithm.
const (
	Minus7 WebAuthnCredentialAlgorithm = -7
)

// Defines values for WebAuthnCredentialType.
const (
	PublicKey WebAuthnCredentialType = "public-key"
)

// Answer 選択したボタンid配列
type Answer = []openapi_types.UUID

// Choice 選択肢ボタン
type Choice struct {
	// Choice 質問文
	Choice string             `json:"choice"`
	Id     openapi_types.UUID `json:"id"`
}

// NewPoll 質問idは存在しない。POST /polls/のボディ。
type NewPoll struct {
	// Deadline deadlineをpostする、またはqStatusがlimitedの時、存在する。回答締め切り時刻。
	Deadline *time.Time `json:"deadline,omitempty"`
	Question []string   `json:"question"`

	// Tags 初期実装では含まない。
	Tags  *[]string `json:"tags,omitempty"`
	Title string    `json:"title"`
	Type  PollType  `json:"type"`
}

// PollBase 質問idは存在しない。
type PollBase struct {
	// Deadline deadlineをpostする、またはqStatusがlimitedの時、存在する。回答締め切り時刻。
	Deadline *time.Time `json:"deadline,omitempty"`

	// Question 質問
	Question Questions `json:"question"`

	// Tags 初期実装では含まない。
	// Deprecated:
	Tags  *[]PollTag `json:"tags,omitempty"`
	Title string     `json:"title"`
	Type  PollType   `json:"type"`
}

// PollComment defines model for PollComment.
type PollComment struct {
	// Content コメント本文
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// PollID defines model for PollID.
type PollID = openapi_types.UUID

// PollResults defines model for PollResults.
type PollResults struct {
	// Count 回答総数
	Count  int      `json:"count"`
	PollId PollID   `json:"pollId"`
	Result []Result `json:"result"`
	Type   PollType `json:"type"`
}

// PollStatus 質問の状態
type PollStatus string

// PollSummaries defines model for PollSummaries.
type PollSummaries = []PollSummary

// PollSummary defines model for PollSummary.
type PollSummary struct {
	CreatedAt time.Time `json:"createdAt"`

	// Deadline deadlineをpostする、またはqStatusがlimitedの時、存在する。回答締め切り時刻。
	Deadline *time.Time `json:"deadline,omitempty"`
	Owner    User       `json:"owner"`
	PollId   PollID     `json:"pollId"`

	// QStatus 質問の状態
	QStatus PollStatus `json:"qStatus"`

	// Question 質問
	Question Questions `json:"question"`

	// Tags 初期実装では含まない。
	// Deprecated:
	Tags  *[]PollTag `json:"tags,omitempty"`
	Title string     `json:"title"`
	Type  PollType   `json:"type"`

	// UserStatus 質問idに対するユーザーの権限
	UserStatus UserStatus `json:"userStatus"`
}

// PollTag defines model for PollTag.
type PollTag struct {
	Id   openapi_types.UUID `json:"id"`
	Name string             `json:"name"`
}

// PollTags defines model for PollTags.
type PollTags = []PollTag

// PollType defines model for PollType.
type PollType string

// PostPollId defines model for PostPollId.
type PostPollId struct {
	// Answer 選択したボタンid配列
	Answer  Answer `json:"answer"`
	Comment string `json:"comment"`
}

// PostTag defines model for PostTag.
type PostTag = string

// PostUser defines model for PostUser.
type PostUser struct {
	// Name アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	Name     UserName     `json:"name"`
	Password UserPassword `json:"password"`
}

// Questions 質問
type Questions = []Choice

// Response defines model for Response.
type Response struct {
	// Answer 選択したボタンid配列
	Answer Answer `json:"answer"`

	// Comment コメント本文
	Comment   *string   `json:"comment,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// Result defines model for Result.
type Result struct {
	// Choice 質問文
	Choice string `json:"choice"`

	// Count その選択肢に回答をした人数
	Count int                `json:"count"`
	Id    openapi_types.UUID `json:"id"`
}

// User defines model for User.
type User struct {
	// Name アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	Name UserName           `json:"name"`
	Uuid openapi_types.UUID `json:"uuid"`
}

// UserName アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
type UserName = string

// UserPassword defines model for UserPassword.
type UserPassword = string

// UserStatus 質問idに対するユーザーの権限
type UserStatus struct {
	// AccessMode only_browable 質問の閲覧　can_answer 解答できる　can_access_details 結果の表示
	AccessMode UserStatusAccessMode `json:"accessMode"`

	// IsOwner オーナーか
	IsOwner bool `json:"isOwner"`
}

// UserStatusAccessMode only_browable 質問の閲覧　can_answer 解答できる　can_access_details 結果の表示
type UserStatusAccessMode string

// WebAuthnAuthenticatorAttachment defines model for WebAuthnAuthenticatorAttachment.
type WebAuthnAuthenticatorAttachment string

// WebAuthnAuthenticatorAttestationType defines model for WebAuthnAuthenticatorAttestationType.
type WebAuthnAuthenticatorAttestationType string

// WebAuthnAuthenticatorResidentKeyRequirement defines model for WebAuthnAuthenticatorResidentKeyRequirement.
type WebAuthnAuthenticatorResidentKeyRequirement string

// WebAuthnAuthenticatorSelectionCriteria defines model for WebAuthnAuthenticatorSelectionCriteria.
type WebAuthnAuthenticatorSelectionCriteria struct {
	AuthenticatorAttachment *WebAuthnAuthenticatorAttachment             `json:"authenticatorAttachment,omitempty"`
	RequireResidentKey      *bool                                        `json:"requireResidentKey,omitempty"`
	ResidentKey             *WebAuthnAuthenticatorResidentKeyRequirement `json:"residentKey,omitempty"`
}

// WebAuthnChallenge defines model for WebAuthnChallenge.
type WebAuthnChallenge = string

// WebAuthnCredential defines model for WebAuthnCredential.
type WebAuthnCredential struct {
	CreatedAt  time.Time          `json:"createdAt"`
	Id         openapi_types.UUID `json:"id"`
	LastUsedAt time.Time          `json:"lastUsedAt"`
	Name       string             `json:"name"`
}

// WebAuthnCredentialAlgorithm defines model for WebAuthnCredentialAlgorithm.
type WebAuthnCredentialAlgorithm int

// WebAuthnCredentialType defines model for WebAuthnCredentialType.
type WebAuthnCredentialType string

// WebAuthnPubKeyCredParam defines model for WebAuthnPubKeyCredParam.
type WebAuthnPubKeyCredParam struct {
	Alg  WebAuthnCredentialAlgorithm `json:"alg"`
	Type WebAuthnCredentialType      `json:"type"`
}

// WebAuthnPublicKeyCredentialCreation defines model for WebAuthnPublicKeyCredentialCreation.
type WebAuthnPublicKeyCredentialCreation struct {
	Id       string                                      `json:"id"`
	RawId    string                                      `json:"rawId"`
	Response WebAuthnPublicKeyCredentialCreationResponse `json:"response"`
	Type     WebAuthnCredentialType                      `json:"type"`
}

// WebAuthnPublicKeyCredentialCreationOptions defines model for WebAuthnPublicKeyCredentialCreationOptions.
type WebAuthnPublicKeyCredentialCreationOptions struct {
	Attestation            WebAuthnAuthenticatorAttestationType   `json:"attestation"`
	AuthenticatorSelection WebAuthnAuthenticatorSelectionCriteria `json:"authenticatorSelection"`
	Challenge              WebAuthnChallenge                      `json:"challenge"`
	PubKeyCredParams       []WebAuthnPubKeyCredParam              `json:"pubKeyCredParams"`
	Rp                     WebAuthnRelyingParty                   `json:"rp"`
	Timeout                int                                    `json:"timeout"`
	User                   WebAuthnUser                           `json:"user"`
}

// WebAuthnPublicKeyCredentialCreationResponse defines model for WebAuthnPublicKeyCredentialCreationResponse.
type WebAuthnPublicKeyCredentialCreationResponse struct {
	AttestationObject string `json:"attestationObject"`
	ClientDataJSON    string `json:"clientDataJSON"`
}

// WebAuthnPublicKeyCredentialRequest defines model for WebAuthnPublicKeyCredentialRequest.
type WebAuthnPublicKeyCredentialRequest struct {
	Id       string                                     `json:"id"`
	RawId    string                                     `json:"rawId"`
	Response WebAuthnPublicKeyCredentialRequestResponse `json:"response"`
	Type     WebAuthnCredentialType                     `json:"type"`
}

// WebAuthnPublicKeyCredentialRequestOptions defines model for WebAuthnPublicKeyCredentialRequestOptions.
type WebAuthnPublicKeyCredentialRequestOptions struct {
	Challenge WebAuthnChallenge `json:"challenge"`
	RpId      *string           `json:"rpId,omitempty"`
	Timeout   *int              `json:"timeout,omitempty"`
}

// WebAuthnPublicKeyCredentialRequestResponse defines model for WebAuthnPublicKeyCredentialRequestResponse.
type WebAuthnPublicKeyCredentialRequestResponse struct {
	AuthenticatorData string  `json:"authenticatorData"`
	ClientDataJSON    string  `json:"clientDataJSON"`
	Signature         string  `json:"signature"`
	UserHandle        *string `json:"userHandle,omitempty"`
}

// WebAuthnRelyingParty defines model for WebAuthnRelyingParty.
type WebAuthnRelyingParty struct {
	Id   *string `json:"id,omitempty"`
	Name string  `json:"name"`
}

// WebAuthnUser defines model for WebAuthnUser.
type WebAuthnUser struct {
	// DisplayName アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	DisplayName UserName           `json:"displayName"`
	Id          openapi_types.UUID `json:"id"`

	// Name アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	Name UserName `json:"name"`
}

// GetPollsParams defines parameters for GetPolls.
type GetPollsParams struct {
	// Limit 最大質問数
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset 質問オフセット
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`

	// Match タイトルの部分一致
	Match *string `form:"match,omitempty" json:"match,omitempty"`
}

// GetPollsPollIDCommentsParams defines parameters for GetPollsPollIDComments.
type GetPollsPollIDCommentsParams struct {
	// Limit 最大コメント取得数
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset オフセット
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// PostPollsJSONRequestBody defines body for PostPolls for application/json ContentType.
type PostPollsJSONRequestBody = NewPoll

// PostPollsPollIDJSONRequestBody defines body for PostPollsPollID for application/json ContentType.
type PostPollsPollIDJSONRequestBody = PostPollId

// PostTagsJSONRequestBody defines body for PostTags for application/json ContentType.
type PostTagsJSONRequestBody = PostTag

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody = PostUser

// PostUsersSigninJSONRequestBody defines body for PostUsersSignin for application/json ContentType.
type PostUsersSigninJSONRequestBody = PostUser

// PostWebauthnAuthenticateFinishJSONRequestBody defines body for PostWebauthnAuthenticateFinish for application/json ContentType.
type PostWebauthnAuthenticateFinishJSONRequestBody = WebAuthnPublicKeyCredentialRequest

// PostWebauthnResisterFinishJSONRequestBody defines body for PostWebauthnResisterFinish for application/json ContentType.
type PostWebauthnResisterFinishJSONRequestBody = WebAuthnPublicKeyCredentialCreation

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /polls)
	GetPolls(ctx echo.Context, params GetPollsParams) error

	// (POST /polls)
	PostPolls(ctx echo.Context) error

	// (DELETE /polls/{pollID})
	DeletePollsPollID(ctx echo.Context, pollID string) error

	// (GET /polls/{pollID})
	GetPollsPollID(ctx echo.Context, pollID string) error

	// (POST /polls/{pollID})
	PostPollsPollID(ctx echo.Context, pollID string) error

	// (POST /polls/{pollID}/close)
	PostPollsClose(ctx echo.Context, pollID string) error

	// (GET /polls/{pollID}/comments)
	GetPollsPollIDComments(ctx echo.Context, pollID string, params GetPollsPollIDCommentsParams) error

	// (GET /polls/{pollID}/results)
	GetPollsPollIDResults(ctx echo.Context, pollID string) error

	// (GET /tags)
	GetTags(ctx echo.Context) error

	// (POST /tags)
	PostTags(ctx echo.Context) error

	// (POST /users)
	PostUsers(ctx echo.Context) error

	// (DELETE /users/me)
	DeleteUsersMe(ctx echo.Context) error

	// (GET /users/me)
	GetUsersMe(ctx echo.Context) error

	// (GET /users/me/answers)
	GetUsersMeAnswers(ctx echo.Context) error

	// (GET /users/me/owners)
	GetUsersMeOwners(ctx echo.Context) error

	// (POST /users/signin)
	PostUsersSignin(ctx echo.Context) error

	// (POST /users/signout)
	PostUsersSignout(ctx echo.Context) error
	// webauthnの認証終了
	// (POST /webauthn/authenticate/finish)
	PostWebauthnAuthenticateFinish(ctx echo.Context) error
	// webauthnの認証開始
	// (POST /webauthn/authenticate/start)
	PostWebauthnAuthenticateStart(ctx echo.Context) error
	// webauthnの登録情報一覧
	// (GET /webauthn/credentials)
	GetWebauthnCredentials(ctx echo.Context) error
	// webauthnの登録情報削除
	// (DELETE /webauthn/credentials/{credentialID})
	DeleteWebauthnCredentials(ctx echo.Context, credentialID openapi_types.UUID) error
	// webauthnの公開鍵登録終了
	// (POST /webauthn/resister/finish)
	PostWebauthnResisterFinish(ctx echo.Context) error
	// webauthnの公開鍵登録開始
	// (POST /webauthn/resister/start)
	PostWebauthnResisterStart(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetPolls converts echo context to params.
func (w *ServerInterfaceWrapper) GetPolls(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPollsParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "match" -------------

	err = runtime.BindQueryParameter("form", true, false, "match", ctx.QueryParams(), &params.Match)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter match: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPolls(ctx, params)
	return err
}

// PostPolls converts echo context to params.
func (w *ServerInterfaceWrapper) PostPolls(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPolls(ctx)
	return err
}

// DeletePollsPollID converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePollsPollID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithOptions("simple", "pollID", ctx.Param("pollID"), &pollID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeletePollsPollID(ctx, pollID)
	return err
}

// GetPollsPollID converts echo context to params.
func (w *ServerInterfaceWrapper) GetPollsPollID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithOptions("simple", "pollID", ctx.Param("pollID"), &pollID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPollsPollID(ctx, pollID)
	return err
}

// PostPollsPollID converts echo context to params.
func (w *ServerInterfaceWrapper) PostPollsPollID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithOptions("simple", "pollID", ctx.Param("pollID"), &pollID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPollsPollID(ctx, pollID)
	return err
}

// PostPollsClose converts echo context to params.
func (w *ServerInterfaceWrapper) PostPollsClose(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithOptions("simple", "pollID", ctx.Param("pollID"), &pollID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPollsClose(ctx, pollID)
	return err
}

// GetPollsPollIDComments converts echo context to params.
func (w *ServerInterfaceWrapper) GetPollsPollIDComments(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithOptions("simple", "pollID", ctx.Param("pollID"), &pollID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPollsPollIDCommentsParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPollsPollIDComments(ctx, pollID, params)
	return err
}

// GetPollsPollIDResults converts echo context to params.
func (w *ServerInterfaceWrapper) GetPollsPollIDResults(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithOptions("simple", "pollID", ctx.Param("pollID"), &pollID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPollsPollIDResults(ctx, pollID)
	return err
}

// GetTags converts echo context to params.
func (w *ServerInterfaceWrapper) GetTags(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTags(ctx)
	return err
}

// PostTags converts echo context to params.
func (w *ServerInterfaceWrapper) PostTags(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostTags(ctx)
	return err
}

// PostUsers converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostUsers(ctx)
	return err
}

// DeleteUsersMe converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteUsersMe(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteUsersMe(ctx)
	return err
}

// GetUsersMe converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersMe(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUsersMe(ctx)
	return err
}

// GetUsersMeAnswers converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersMeAnswers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUsersMeAnswers(ctx)
	return err
}

// GetUsersMeOwners converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersMeOwners(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUsersMeOwners(ctx)
	return err
}

// PostUsersSignin converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsersSignin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostUsersSignin(ctx)
	return err
}

// PostUsersSignout converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsersSignout(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostUsersSignout(ctx)
	return err
}

// PostWebauthnAuthenticateFinish converts echo context to params.
func (w *ServerInterfaceWrapper) PostWebauthnAuthenticateFinish(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostWebauthnAuthenticateFinish(ctx)
	return err
}

// PostWebauthnAuthenticateStart converts echo context to params.
func (w *ServerInterfaceWrapper) PostWebauthnAuthenticateStart(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostWebauthnAuthenticateStart(ctx)
	return err
}

// GetWebauthnCredentials converts echo context to params.
func (w *ServerInterfaceWrapper) GetWebauthnCredentials(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetWebauthnCredentials(ctx)
	return err
}

// DeleteWebauthnCredentials converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteWebauthnCredentials(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "credentialID" -------------
	var credentialID openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "credentialID", ctx.Param("credentialID"), &credentialID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter credentialID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteWebauthnCredentials(ctx, credentialID)
	return err
}

// PostWebauthnResisterFinish converts echo context to params.
func (w *ServerInterfaceWrapper) PostWebauthnResisterFinish(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostWebauthnResisterFinish(ctx)
	return err
}

// PostWebauthnResisterStart converts echo context to params.
func (w *ServerInterfaceWrapper) PostWebauthnResisterStart(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostWebauthnResisterStart(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/polls", wrapper.GetPolls)
	router.POST(baseURL+"/polls", wrapper.PostPolls)
	router.DELETE(baseURL+"/polls/:pollID", wrapper.DeletePollsPollID)
	router.GET(baseURL+"/polls/:pollID", wrapper.GetPollsPollID)
	router.POST(baseURL+"/polls/:pollID", wrapper.PostPollsPollID)
	router.POST(baseURL+"/polls/:pollID/close", wrapper.PostPollsClose)
	router.GET(baseURL+"/polls/:pollID/comments", wrapper.GetPollsPollIDComments)
	router.GET(baseURL+"/polls/:pollID/results", wrapper.GetPollsPollIDResults)
	router.GET(baseURL+"/tags", wrapper.GetTags)
	router.POST(baseURL+"/tags", wrapper.PostTags)
	router.POST(baseURL+"/users", wrapper.PostUsers)
	router.DELETE(baseURL+"/users/me", wrapper.DeleteUsersMe)
	router.GET(baseURL+"/users/me", wrapper.GetUsersMe)
	router.GET(baseURL+"/users/me/answers", wrapper.GetUsersMeAnswers)
	router.GET(baseURL+"/users/me/owners", wrapper.GetUsersMeOwners)
	router.POST(baseURL+"/users/signin", wrapper.PostUsersSignin)
	router.POST(baseURL+"/users/signout", wrapper.PostUsersSignout)
	router.POST(baseURL+"/webauthn/authenticate/finish", wrapper.PostWebauthnAuthenticateFinish)
	router.POST(baseURL+"/webauthn/authenticate/start", wrapper.PostWebauthnAuthenticateStart)
	router.GET(baseURL+"/webauthn/credentials", wrapper.GetWebauthnCredentials)
	router.DELETE(baseURL+"/webauthn/credentials/:credentialID", wrapper.DeleteWebauthnCredentials)
	router.POST(baseURL+"/webauthn/resister/finish", wrapper.PostWebauthnResisterFinish)
	router.POST(baseURL+"/webauthn/resister/start", wrapper.PostWebauthnResisterStart)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xcf3PTRvp/K8x+v3/a2E4CBf+Xhmsv1yvJJTCdKZNhNtbGUU+WjLQmpBnPRHIhOQID",
	"w5VwzIXSFi4JoTW0lB7QXnkxix3nXdzsrn5ZWskrSEp7/zB2vHr2+fl59nn2EUugYtTqho50bIHyErAq",
	"86gG2cdR3VpAJv2kIKtiqnWsGjoogz37WffyF8S+Rey7pLVBnJek9URV9i5e7azeAjmgYlRjBOYMswYx",
	"KINGQ1VADuDFOgJlYGFT1augmQM1VR/ni0s5gFWs0Z/dbf3l0DThIl09Nm+oFZTET8/52mcG5EDdNOrI",
	"xCpinFQSHu092encvNZdXwE5UIMX/oz0Kp4H5eFikTHnfS8JeFcVCRGbOWCicw3VRAoonwFsicvLTCCx",
	"K5j/tDH7CapgkAMX8ugCrNU1LoT7OV8KS+RtxRkq0h0v5FUdI1OHGijPQc1CzRw4iRYmDU1L0oCqEPtR",
	"59t/dDa2mV13iP0ZWXYmJ6ZPHSrUDU2zCsRuU/22Vohzjyw7MRUrCCqaqguU7P1CnBt1w8LEvk2cNbJs",
	"E/sX6kH2o3PTGOKGRewrmlpTMVKI3e7edsiy7bHEn3A6//xi99vPd//9JXHszuoKcS53bzud1Z84P74t",
	"FIhRHqs1JPK5cw1kccaWAk8d4JkRR8SwasXF7Kze6W7c7bTv9u5dJPYWVej1h0xGV5vh0IhtGNuD+8ZS",
	"Vrfkf1gC/2+iOVAG/1cIwrvgxnaBesIpui7qn3xTl0hIVSFf9RxpgLNSynTdu9BCGZ3uf8Kz0tT/F3ed",
	"1e9LdRNVIKaGwGYD5d7EuwYaH0o4ucABSyNv3wF9r5L0wDGjVkM6pixFcoKhY/eHflUT5wlpfUVaT0hr",
	"tbvxDc8OMTkrJqLWGsV9aSDFQSKiBs/nfFYiYnqsZ0sMnliAOJ8R+x6xvzuv1tiHW8ShooQ4B0PF0vF8",
	"8Vh+aORUaaRcOlIeGvoYND3djZ8Q5rgQi+MnYoJG7SBOSPThKWQ1NH7siJqmITKMF6TXujcfB9tS2lVk",
	"UpvQTDWuyHjf+AnA7EEZ6MsEaQ9yfoVo/SZO73Ltez2X3ucu4hSe0qJO4eqUA14S5BK7vXv5x+7FNZAD",
	"SG/U6PZGHemIbu9CJMgBo4EVBkUzIuOKjTndqNWg6RpQGon4U2HAiRCL6zr8WHkJQE2bmAPlM1EfyuYL",
	"YpPMxJU8mBwDJ7ow4tNZ8SIHjAWdn7/Ttjxtcd8/F1h+oNL5ymYONCxkyj13OliZhmQeFx7zfVvE1TkT",
	"TXPBiZOdN78hznPSukOx2HlOlp33/+CfR8v0+CBYAkSetJiUHWgmjMGP1NE+B3RYiybII4Pyo6ggYHQi",
	"UU75GpTgkmPxlHuoyHogiHCQFH2nXLTzEMSEimrEJOBgJgcfFp70w7XfFtAvQtNkcGtGmpaDdJ+ueuiV",
	"md4Tffz7DAlh1sKu34Qf6DNZ4CX0JxajMdE8/xkUdifpOprdoGUtGKYi88yktzYqNts0RCsiNeNUIHNw",
	"Zk3ILLLHT7fWTSr+g30irieI3ylk1Q2d1xf75jRv4ygYuKJPSJB5wJR/WglynpSqJY9XxL5D7HbQTbEf",
	"8gMXcW7wZs+rFy/EJ69oQmD0JdB+9+n17hcbIcR2RczY+XDlKfo9kGYO7E/IMejP3OZxF/m4Ln945/yx",
	"nOkGKucAFI8fmUVHS0fy8BhU8iOV4eH8MXS8kq+8c7Q0NDc7d7RyfA4kNX98iQT+/TVxHhJnk7t45/pV",
	"suzQLYm9tdv+avf6JWLfJM4VWiTbV2hl3dokrZ+J8yNp/dzbXN91nhO73dt8wj5s7a1c7d1fIc6N3vZ3",
	"9EFWm/a31kpH+zLkCIUiTBkGZXCmmD8O85+O5j8+OxNyCp99Qej1YV1qLj4muZNPLWG39AM21dzDzqNf",
	"3MZCSFvEbne3H+zdvh7rb8BKBVnWh4YiMJCha4tnZ01jAc5q6JB/iN9b/763uUWWlytQP8vh41Bv6x4N",
	"VnuL2FdZU4P/yIifVRCGqmYd4hFHbfbV9u79F+EywNvJolvRqPJJe1/6SMVrg2YOqNaEd2aNOtoO1UHr",
	"MtPEWqDcWcPQENTjJyOXUi6sHmEwxf39IzQ72sDzOv0H6VitQGyYoxjDyrwH857UdQ1iGtwMfA3Lyvt/",
	"CCXGQfQEikh6BFkYUpVEj1C6oVOdq7qimlww94MEG300ZXmZQpaqIB1/gBanuNqjqvGtQf0VzSGTf1ZU",
	"q2I0TFhFykDmEjaR5XEaaahCJRszVYxMFQoyfbKJ0zB+kEUDZwyJwL16DrIkzHt0USdmRXt4fWYmElTW",
	"bA5QdVxXghOE9+TYPNQ0pFdRX26bhRY6OgLiGwXLUyw3ZiLKtwo1QT8ne+0rWYNpkJ1ZM5He17ot11f+",
	"htgRxEZIRWnG8VeNalXDVPF8LRSW+XdSCQePiPpj8fVRHKo3ZjW1kv8rWkwXYBDWTDZmP0CLdP0kNGFN",
	"ELhaVTY+RNJJ9twSGI+1nbk0lCeB1FFZUmw3ydTnruYbjlH3cG8FUtsMQfhFlWrChXHptaGiSEYzKRz7",
	"9dX+KpvFDhcpxK67hVj9iTp9PVNM1P1SNuKUQTJ93QzSl4ybuf4E5YP0a1GPQzytPcNQLmUd/4FmjoZ7",
	"2LXlm0ZJsSHojZt1WWJTSFtU9eokNLHbEa4ho4H78u7RYrFYFGFbwxpc8Xsb8Z5pxC/NututBGGtCnQU",
	"MJZo3lyfK2Vzas87X8+3U5oiAUcTnKIcolQ0Fen4BMTwT9MTJ6UeivYD+inkBKxkU9FUBDUkdUQPVMjC",
	"vx8kdhn+NYDYTYE+03IG8TT6WnZIhOE3wzSzHjFRw1RF9smALlF/9rfLpKXXC+yoD6SXQDTIDi6uc8BS",
	"qzrEDRNJrqd4+keoKxraH+CIyRpmKZM1ZDCkLyENvCYSu9k+VBvRCyIheyliiHuhimrVNbh4MmNLNOPl",
	"mBzZlPIqzKZAB+ILC0pQ1ecMbyAC8mzntldHJ8cPTTfqdcOkyNUwNVAG8xjXy4XCwsLCYbcje7hi1AqW",
	"u6wZbV2PTo6HKmX+7TwyLf5r6XDxcJFd4daRDusqKIPhw8XDQ+zmBc8z/fN7TPqpigTd+O7lm7v/2iat",
	"HXa/uUqcG51r651fbvkjS6S1Qlo3ifOAtB7SBWwoqLv+mPXqP9v78hKxH9LF9nNib8YfB4w5k+VTCpjg",
	"fcTuvCzGowlrCCPTYpcMEcY2ljv3t9zJSXYdoNI/n2sgds3q6pgNE9DoZLYOXceFQDVhRsHZYXL9RFot",
	"0lpNoG/MzVko6wbEeUmc+1RZrYfEbu+1tjurl149W+6t/JCwTQ3iyrxoFz9UZ4KkyYw5VCxGpotgva5R",
	"xFINvfCJxQ//AT25CQkas8yr+wWa+IB3hth9Lx9bADNsEMZKdinnxqv/bHRXr3NXiPmBd/lpAR6UyMLv",
	"GsrivgnlTRA2m81mTHelA9DdokhzY7xrQ2N0hFus//d3oXLIO92wNaX4mtM6zUeGqX6KFJEdmjk3yAtL",
	"bKDkRJPT0BBOHIqkgfq3y3u37ydY5wR7mtnHn7+K6HAkTpv6iZQQdNFwfNF7hjmrKgrS+QrBDicNfOg9",
	"o6ErYocUQtz4CX5Z0nm5waV1HdTe5leR7vWj3e49eLL7w+NExEpSxNtxJikdZ9dgBJMZWNFcEmBVPVCD",
	"l0p5lzoVvMRQkWiaZ75dkoFjMszIfsNHaDZDiCAy6Jt4/3oeag0UniA4I3fzOhOaIQjdS0sNWsqJHdSA",
	"4ixwYJ4XB7FCRTPcQuQ34ZPODcbQoHQ2xrjed7TcD31y10k+C8YkD8+l8PduMhzweHiOeZtKHffCG/J9",
	"9vfot3+HvpkDQYQz4cFqNk99i9hrnbU7xL5LnL/LBHvQonzTKY+ZftiQnvLz5sljTdpfC1L8UbsDy2eC",
	"6DKDKXNhcAWT0Xw2wrnRe/k5sW8PiJ9gDvtAElAw0+QNMoN3SkNDw0UF5pWjx4bzI8eHYf4YHFHyFThS",
	"QiMVNDs0PByMjTOXHTAqJUeyGYxd8EnPDGkrPLH+q7mZNzh/oF7mvbgj9Kmg7rbbrOh8HMFnkXO5E7cH",
	"Wk2yPQYXkhhWU+pIX6Lu+uPe5rVwNSnKO+5wrHWAB0I2v5xYT75B5RfTCTU9BW73HQOxevqG6rh6hFo5",
	"zQgdnFrce66kU3KqC7BkFZK34M0PiqvX6CQhr2GXHQ9dd4j9gNjtbuti58vviP2o215LcBde4TLVfIjA",
	"63OeUHT2VnY6q5fYqwPBcJ7L1uAYTWdrXwznGS2X2T4FXrUkwxIv3LhJXj1b7m1upcg46hL7DfW3RDKz",
	"90ySRebRJynyBKf1W5bYUqu6qqeAT+tbhs33aRwmV0RM4GlO63eCP1Ry974uCXefunJTLNpkPfHBCuD3",
	"+G/C3gKapUcVvRC6nUKFOVVXrflkbr0LjN7O1d72z7tPnVcvLgk5/cglH5oFQe9x4gdjOYn73iw2lW2e",
	"Wt47fcBTKLHbEe142vdWpFrAwtDEsgbYW1/rbK1JG2Ca0T5AoJC/Tk4GkRSN+uIO0mjF3zYZYz1ed2//",
	"tHfle55Lk7HWU+dYiPIbKjLTyFJoClOyEE5SpVDeDAotLAVfBlwOiFTMz1gJ5ycpLY8kwtxAef3NBfLK",
	"lF1hyVOLr0Gv3sz0addElmphZErDL5dJAn6nXMpvC3r9mUfZOmcKVRnDGUododU7F7/ZW1/bu/o0oqpB",
	"ju6bQhKHOXkJHPYs8TYxODqslxE5IjpNRWNKA5nnvVji4woFWFd5Y4Y/0ddUZK/8ud99SqG/sU506Dtm",
	"b//6X03vvznw/+L17Zozzf8GAAD//+KJqXswSgAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
