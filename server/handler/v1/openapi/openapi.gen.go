// Package Openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
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
type Answer []int

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
	Count  *int    `json:"count,omitempty"`
	PollId *PollID `json:"pollId,omitempty"`

	// 結果
	Result *Results  `json:"result,omitempty"`
	Type   *PollType `json:"type,omitempty"`
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
	Answer  *Answer `json:"answer,omitempty"`
	Comment *string `json:"comment,omitempty"`
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

// 結果
type Results []struct {
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
	"W1bXSrX/BAAA///ddWVauy8AAA==",
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
