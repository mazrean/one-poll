//go:build wireinject
// +build wireinject

package main

import (
	v1Handler "github.com/cs-sysimpl/suzukake/handler/v1"
	"github.com/cs-sysimpl/suzukake/pkg/common"
	"github.com/cs-sysimpl/suzukake/repository"
	"github.com/cs-sysimpl/suzukake/repository/gorm2"
	"github.com/cs-sysimpl/suzukake/service"
	v1Service "github.com/cs-sysimpl/suzukake/service/v1"
	"github.com/google/wire"
)

type Config struct {
	IsProduction  common.IsProduction
	SessionKey    common.SessionKey
	SessionSecret common.SessionSecret
}

var (
	isProductionField  = wire.FieldsOf(new(*Config), "IsProduction")
	sessionKeyField    = wire.FieldsOf(new(*Config), "SessionKey")
	sessionSecretField = wire.FieldsOf(new(*Config), "SessionSecret")
)

var (
	dbBind                 = wire.Bind(new(repository.DB), new(*gorm2.DB))
	userRepositoryBind     = wire.Bind(new(repository.User), new(*gorm2.User))
	pollRepositoryBind     = wire.Bind(new(repository.Poll), new(*gorm2.Poll))
	choiceRepositoryBind   = wire.Bind(new(repository.Choice), new(*gorm2.Choice))
	tagRepositoryBind      = wire.Bind(new(repository.Tag), new(*gorm2.Tag))
	responseRepositoryBind = wire.Bind(new(repository.Response), new(*gorm2.Response))
	commentRepositoryBind  = wire.Bind(new(repository.Comment), new(*gorm2.Comment))

	authorizationServiceBind = wire.Bind(new(service.Authorization), new(*v1Service.Authorization))
	pollServiceBind          = wire.Bind(new(service.Poll), new(*v1Service.Poll))
	tagServiceBind           = wire.Bind(new(service.Tag), new(*v1Service.Tag))
	commentServiceBind       = wire.Bind(new(service.Comment), new(*v1Service.Comment))
	responseServiceBind      = wire.Bind(new(service.Response), new(*v1Service.Response))
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

		dbBind,
		userRepositoryBind,
		pollRepositoryBind,
		choiceRepositoryBind,
		tagRepositoryBind,
		responseRepositoryBind,
		commentRepositoryBind,
		gorm2.NewDB,
		gorm2.NewUser,
		gorm2.NewPoll,
		gorm2.NewChoice,
		gorm2.NewTag,
		gorm2.NewResponse,
		gorm2.NewComment,

		authorizationServiceBind,
		pollServiceBind,
		tagServiceBind,
		commentServiceBind,
		responseServiceBind,
		v1Service.NewAuthorization,
		v1Service.NewPoll,
		v1Service.NewTag,
		v1Service.NewComment,
		v1Service.NewResponse,
		v1Service.NewPollAuthority,

		v1Handler.NewAPI,
		v1Handler.NewSession,
		//v1Handler.NewContext,
		v1Handler.NewChecker,
		v1Handler.NewUser,
		v1Handler.NewPoll,
		v1Handler.NewTag,
		v1Handler.NewComment,
		v1Handler.NewResponse,

		NewService,
	)
	return nil, nil
}
