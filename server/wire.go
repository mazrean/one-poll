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
	dbBind             = wire.Bind(new(repository.DB), new(*gorm2.DB))
	userRepositoryBind = wire.Bind(new(repository.User), new(*gorm2.User))

	authorizationServiceBind = wire.Bind(new(service.Authorization), new(*v1Service.Authorization))
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
		gorm2.NewDB,
		gorm2.NewUser,

		authorizationServiceBind,
		v1Service.NewAuthorization,

		v1Handler.NewAPI,
		v1Handler.NewSession,
		//v1Handler.NewContext,
		v1Handler.NewChecker,
		v1Handler.NewUser,

		NewService,
	)
	return nil, nil
}
