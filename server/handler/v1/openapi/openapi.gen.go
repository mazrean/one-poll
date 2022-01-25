// Package Openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package Openapi

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

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Defines values for PollStatus.
const (
	PollStatusLimited PollStatus = "limited"

	PollStatusOpened PollStatus = "opened"

	PollStatusOutdated PollStatus = "outdated"
)

// Defines values for PollType.
const (
	PollTypeRadio PollType = "radio"
)

// Defines values for UserStatusAccessMode.
const (
	UserStatusAccessModeCanAccessDetails UserStatusAccessMode = "can_access_details"

	UserStatusAccessModeCanAnswer UserStatusAccessMode = "can_answer"

	UserStatusAccessModeOnlyBrowsable UserStatusAccessMode = "only_browsable"
)

// 選択したボタンid配列
type Answer []string

// 選択肢ボタン
type Choice struct {
	// 質問文
	Choice string `json:"choice"`
	Id     string `json:"id"`
}

// 質問idは存在しない。POST /polls/のボディ。
type NewPoll struct {
	// deadlineをpostする、またはqStatusがlimitedの時、存在する。回答締め切り時刻。
	Deadline *time.Time `json:"deadline,omitempty"`
	Question []string   `json:"question"`

	// 初期実装では含まない。
	Tags  *[]string `json:"tags,omitempty"`
	Title string    `json:"title"`
	Type  PollType  `json:"type"`
}

// 質問idは存在しない。
type PollBase struct {
	// deadlineをpostする、またはqStatusがlimitedの時、存在する。回答締め切り時刻。
	Deadline *time.Time `json:"deadline,omitempty"`

	// 質問
	Question Questions `json:"question"`

	// 初期実装では含まない。
	Tags  *[]PollTag `json:"tags,omitempty"`
	Title string     `json:"title"`
	Type  PollType   `json:"type"`
}

// PollComment defines model for PollComment.
type PollComment struct {
	// コメント本文
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

// PollID defines model for PollID.
type PollID string

// PollResults defines model for PollResults.
type PollResults struct {
	// 回答総数
	Count  int      `json:"count"`
	PollId PollID   `json:"pollId"`
	Result []Result `json:"result"`
	Type   PollType `json:"type"`
}

// 質問の状態
type PollStatus string

// PollSummaries defines model for PollSummaries.
type PollSummaries []PollSummary

// PollSummary defines model for PollSummary.
type PollSummary struct {
	// Embedded fields due to inline allOf schema
	PollId PollID `json:"pollId"`
	// Embedded struct due to allOf(#/components/schemas/PollBase)
	PollBase `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	CreatedAt time.Time `json:"createdAt"`
	Owner     User      `json:"owner"`

	// 質問の状態
	QStatus PollStatus `json:"qStatus"`

	// 質問idに対するユーザーの権限
	UserStatus UserStatus `json:"userStatus"`
}

// PollTag defines model for PollTag.
type PollTag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// PollTags defines model for PollTags.
type PollTags []PollTag

// PollType defines model for PollType.
type PollType string

// PostPollId defines model for PostPollId.
type PostPollId struct {
	// 選択したボタンid配列
	Answer  Answer `json:"answer"`
	Comment string `json:"comment"`
}

// PostTag defines model for PostTag.
type PostTag string

// PostUser defines model for PostUser.
type PostUser struct {
	// アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	Name     UserName     `json:"name"`
	Password UserPassword `json:"password"`
}

// 質問
type Questions []Choice

// Response defines model for Response.
type Response struct {
	// 選択したボタンid配列
	Answer Answer `json:"answer"`

	// コメント本文
	Comment   *string   `json:"comment,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// Result defines model for Result.
type Result struct {
	// Embedded struct due to allOf(#/components/schemas/Choice)
	Choice `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// その選択肢に回答をした人数
	Count int `json:"count"`
}

// User defines model for User.
type User struct {
	// アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	Name UserName `json:"name"`
	Uuid string   `json:"uuid"`
}

// アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
type UserName string

// UserPassword defines model for UserPassword.
type UserPassword string

// 質問idに対するユーザーの権限
type UserStatus struct {
	// only_browable 質問の閲覧　can_answer 解答できる　can_access_details 結果の表示
	AccessMode UserStatusAccessMode `json:"accessMode"`

	// オーナーか
	IsOwner bool `json:"isOwner"`
}

// only_browable 質問の閲覧　can_answer 解答できる　can_access_details 結果の表示
type UserStatusAccessMode string

// GetPollsParams defines parameters for GetPolls.
type GetPollsParams struct {
	// 最大質問数
	Limit *int `json:"limit,omitempty"`

	// 質問オフセット
	Offset *int `json:"offset,omitempty"`

	// タイトルの部分一致
	Match *string `json:"match,omitempty"`
}

// PostPollsJSONBody defines parameters for PostPolls.
type PostPollsJSONBody NewPoll

// PostPollsPollIDJSONBody defines parameters for PostPollsPollID.
type PostPollsPollIDJSONBody PostPollId

// GetPollsPollIDCommentsParams defines parameters for GetPollsPollIDComments.
type GetPollsPollIDCommentsParams struct {
	// 最大コメント取得数
	Limit *int `json:"limit,omitempty"`

	// オフセット
	Offset *int `json:"offset,omitempty"`
}

// PostTagsJSONBody defines parameters for PostTags.
type PostTagsJSONBody PostTag

// PostUsersJSONBody defines parameters for PostUsers.
type PostUsersJSONBody PostUser

// PostUsersSigninJSONBody defines parameters for PostUsersSignin.
type PostUsersSigninJSONBody PostUser

// PostPollsJSONRequestBody defines body for PostPolls for application/json ContentType.
type PostPollsJSONRequestBody PostPollsJSONBody

// PostPollsPollIDJSONRequestBody defines body for PostPollsPollID for application/json ContentType.
type PostPollsPollIDJSONRequestBody PostPollsPollIDJSONBody

// PostTagsJSONRequestBody defines body for PostTags for application/json ContentType.
type PostTagsJSONRequestBody PostTagsJSONBody

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody PostUsersJSONBody

// PostUsersSigninJSONRequestBody defines body for PostUsersSignin for application/json ContentType.
type PostUsersSigninJSONRequestBody PostUsersSigninJSONBody

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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPolls(ctx, params)
	return err
}

// PostPolls converts echo context to params.
func (w *ServerInterfaceWrapper) PostPolls(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostPolls(ctx)
	return err
}

// DeletePollsPollID converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePollsPollID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "pollID", runtime.ParamLocationPath, ctx.Param("pollID"), &pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeletePollsPollID(ctx, pollID)
	return err
}

// GetPollsPollID converts echo context to params.
func (w *ServerInterfaceWrapper) GetPollsPollID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "pollID", runtime.ParamLocationPath, ctx.Param("pollID"), &pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPollsPollID(ctx, pollID)
	return err
}

// PostPollsPollID converts echo context to params.
func (w *ServerInterfaceWrapper) PostPollsPollID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "pollID", runtime.ParamLocationPath, ctx.Param("pollID"), &pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostPollsPollID(ctx, pollID)
	return err
}

// PostPollsClose converts echo context to params.
func (w *ServerInterfaceWrapper) PostPollsClose(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "pollID", runtime.ParamLocationPath, ctx.Param("pollID"), &pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostPollsClose(ctx, pollID)
	return err
}

// GetPollsPollIDComments converts echo context to params.
func (w *ServerInterfaceWrapper) GetPollsPollIDComments(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "pollID", runtime.ParamLocationPath, ctx.Param("pollID"), &pollID)
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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPollsPollIDComments(ctx, pollID, params)
	return err
}

// GetPollsPollIDResults converts echo context to params.
func (w *ServerInterfaceWrapper) GetPollsPollIDResults(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "pollID" -------------
	var pollID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "pollID", runtime.ParamLocationPath, ctx.Param("pollID"), &pollID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pollID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPollsPollIDResults(ctx, pollID)
	return err
}

// GetTags converts echo context to params.
func (w *ServerInterfaceWrapper) GetTags(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTags(ctx)
	return err
}

// PostTags converts echo context to params.
func (w *ServerInterfaceWrapper) PostTags(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostTags(ctx)
	return err
}

// PostUsers converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostUsers(ctx)
	return err
}

// DeleteUsersMe converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteUsersMe(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteUsersMe(ctx)
	return err
}

// GetUsersMe converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersMe(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUsersMe(ctx)
	return err
}

// GetUsersMeAnswers converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersMeAnswers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUsersMeAnswers(ctx)
	return err
}

// GetUsersMeOwners converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersMeOwners(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUsersMeOwners(ctx)
	return err
}

// PostUsersSignin converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsersSignin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostUsersSignin(ctx)
	return err
}

// PostUsersSignout converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsersSignout(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostUsersSignout(ctx)
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

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"AAD//6k+paZFMAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
