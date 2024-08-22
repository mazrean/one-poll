//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mazrean/one-poll/domain"
	v1Handler "github.com/mazrean/one-poll/handler/v1"
	"github.com/mazrean/one-poll/pkg/common"
	"github.com/mazrean/one-poll/repository"
	"github.com/mazrean/one-poll/repository/gorm2"
	"github.com/mazrean/one-poll/service"
	v1Service "github.com/mazrean/one-poll/service/v1"
)

type Config struct {
	IsProduction  common.IsProduction
	SessionKey    common.SessionKey
	SessionSecret common.SessionSecret
	RelyingParty  *domain.WebAuthnRelyingParty
}

var (
	isProductionField  = wire.FieldsOf(new(*Config), "IsProduction")
	sessionKeyField    = wire.FieldsOf(new(*Config), "SessionKey")
	sessionSecretField = wire.FieldsOf(new(*Config), "SessionSecret")
	relyingPartyField  = wire.FieldsOf(new(*Config), "RelyingParty")
)

var (
	dbBind                           = wire.Bind(new(repository.DB), new(*gorm2.DB))
	userRepositoryBind               = wire.Bind(new(repository.User), new(*gorm2.User))
	pollRepositoryBind               = wire.Bind(new(repository.Poll), new(*gorm2.Poll))
	choiceRepositoryBind             = wire.Bind(new(repository.Choice), new(*gorm2.Choice))
	tagRepositoryBind                = wire.Bind(new(repository.Tag), new(*gorm2.Tag))
	responseRepositoryBind           = wire.Bind(new(repository.Response), new(*gorm2.Response))
	commentRepositoryBind            = wire.Bind(new(repository.Comment), new(*gorm2.Comment))
	webauthnCredentialRepositoryBind = wire.Bind(new(repository.WebAuthnCredential), new(*gorm2.WebAuthnCredential))

	authorizationServiceBind = wire.Bind(new(service.Authorization), new(*v1Service.Authorization))
	pollServiceBind          = wire.Bind(new(service.Poll), new(*v1Service.Poll))
	tagServiceBind           = wire.Bind(new(service.Tag), new(*v1Service.Tag))
	commentServiceBind       = wire.Bind(new(service.Comment), new(*v1Service.Comment))
	responseServiceBind      = wire.Bind(new(service.Response), new(*v1Service.Response))
	webauthnServiceBind      = wire.Bind(new(service.WebAuthn), new(*v1Service.WebAuthn))
)

type Service struct {
	*v1Handler.API
}

func NewService(api *v1Handler.API) *Service {
	return &Service{
		API: api,
	}
}

func InjectService(config *Config) (*Service, error) {
	wire.Build(
		isProductionField,
		sessionKeyField,
		sessionSecretField,
		relyingPartyField,

		dbBind,
		userRepositoryBind,
		pollRepositoryBind,
		choiceRepositoryBind,
		tagRepositoryBind,
		responseRepositoryBind,
		commentRepositoryBind,
		webauthnCredentialRepositoryBind,
		gorm2.NewDB,
		gorm2.NewUser,
		gorm2.NewPoll,
		gorm2.NewChoice,
		gorm2.NewTag,
		gorm2.NewResponse,
		gorm2.NewComment,
		gorm2.NewWebAuthnCredential,

		authorizationServiceBind,
		pollServiceBind,
		tagServiceBind,
		commentServiceBind,
		responseServiceBind,
		webauthnServiceBind,
		v1Service.NewAuthorization,
		v1Service.NewPoll,
		v1Service.NewTag,
		v1Service.NewComment,
		v1Service.NewResponse,
		v1Service.NewPollAuthority,
		v1Service.NewWebAuthn,

		v1Handler.NewAPI,
		v1Handler.NewSession,
		//v1Handler.NewContext,
		v1Handler.NewChecker,
		v1Handler.NewUser,
		v1Handler.NewPoll,
		v1Handler.NewTag,
		v1Handler.NewComment,
		v1Handler.NewResponse,
		v1Handler.NewWebAuthn,

		NewService,
	)
	return nil, nil
}
