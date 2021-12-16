//go:build wireinject
// +build wireinject

package main

import (
	v1Handler "github.com/cs-sysimpl/suzukake/handler/v1"
	"github.com/cs-sysimpl/suzukake/pkg/common"
	"github.com/cs-sysimpl/suzukake/repository"
	"github.com/cs-sysimpl/suzukake/repository/gorm2"
	"github.com/google/wire"
	//"github.com/cs-sysimpl/suzukake/service"
	//v1Service "github.com/cs-sysimpl/suzukake/service/v1"
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
	dbBind = wire.Bind(new(repository.DB), new(*gorm2.DB))
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
		//isProductionField,
		//sessionKeyField,
		//sessionSecretField,
		//dbBind,
		//gorm2.NewDB,
		v1Handler.NewAPI,
		//v1Handler.NewSession,
		//v1Handler.NewContext,
		v1Handler.NewChecker,
		NewService,
	)
	return nil, nil
}
