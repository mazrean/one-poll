// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
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

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
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

// 選択したボタンid配列
type Answer = []openapi_types.UUID

// 選択肢ボタン
type Choice struct {
	// 質問文
	Choice string             `json:"choice"`
	Id     openapi_types.UUID `json:"id"`
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
}

// PollID defines model for PollID.
type PollID = openapi_types.UUID

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
type PollSummaries = []PollSummary

// PollSummary defines model for PollSummary.
type PollSummary struct {
	CreatedAt time.Time `json:"createdAt"`

	// deadlineをpostする、またはqStatusがlimitedの時、存在する。回答締め切り時刻。
	Deadline *time.Time `json:"deadline,omitempty"`
	Owner    User       `json:"owner"`
	PollId   PollID     `json:"pollId"`

	// 質問の状態
	QStatus PollStatus `json:"qStatus"`

	// 質問
	Question Questions `json:"question"`

	// 初期実装では含まない。
	Tags  *[]PollTag `json:"tags,omitempty"`
	Title string     `json:"title"`
	Type  PollType   `json:"type"`

	// 質問idに対するユーザーの権限
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
	// 選択したボタンid配列
	Answer  Answer `json:"answer"`
	Comment string `json:"comment"`
}

// PostTag defines model for PostTag.
type PostTag = string

// PostUser defines model for PostUser.
type PostUser struct {
	// アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	Name     UserName     `json:"name"`
	Password UserPassword `json:"password"`
}

// 質問
type Questions = []Choice

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
	// 質問文
	Choice string `json:"choice"`

	// その選択肢に回答をした人数
	Count int                `json:"count"`
	Id    openapi_types.UUID `json:"id"`
}

// User defines model for User.
type User struct {
	// アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
	Name UserName           `json:"name"`
	Uuid openapi_types.UUID `json:"uuid"`
}

// アカウント名。uuidで管理されるが、ユーザー視点の観点で重複を許さない
type UserName = string

// UserPassword defines model for UserPassword.
type UserPassword = string

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
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// 質問オフセット
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`

	// タイトルの部分一致
	Match *string `form:"match,omitempty" json:"match,omitempty"`
}

// PostPollsJSONBody defines parameters for PostPolls.
type PostPollsJSONBody = NewPoll

// PostPollsPollIDJSONBody defines parameters for PostPollsPollID.
type PostPollsPollIDJSONBody = PostPollId

// GetPollsPollIDCommentsParams defines parameters for GetPollsPollIDComments.
type GetPollsPollIDCommentsParams struct {
	// 最大コメント取得数
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// オフセット
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// PostTagsJSONBody defines parameters for PostTags.
type PostTagsJSONBody = PostTag

// PostUsersJSONBody defines parameters for PostUsers.
type PostUsersJSONBody = PostUser

// PostUsersSigninJSONBody defines parameters for PostUsersSignin.
type PostUsersSigninJSONBody = PostUser

// PostPollsJSONRequestBody defines body for PostPolls for application/json ContentType.
type PostPollsJSONRequestBody = PostPollsJSONBody

// PostPollsPollIDJSONRequestBody defines body for PostPollsPollID for application/json ContentType.
type PostPollsPollIDJSONRequestBody = PostPollsPollIDJSONBody

// PostTagsJSONRequestBody defines body for PostTags for application/json ContentType.
type PostTagsJSONRequestBody = PostTagsJSONBody

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody = PostUsersJSONBody

// PostUsersSigninJSONRequestBody defines body for PostUsersSignin for application/json ContentType.
type PostUsersSigninJSONRequestBody = PostUsersSigninJSONBody

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
