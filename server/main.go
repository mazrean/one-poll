//go:generate go run github.com/google/wire/cmd/wire@latest

package main

import (
	"fmt"
	"os"

	"github.com/cs-sysimpl/suzukake/pkg/common"
)

func main() {
	env := os.Getenv("APP_ENV")
	isProduction := env == "production"

	secret, ok := os.LookupEnv("SESSION_SECRET")
	if !ok {
		panic("SESSION_SECRET is not set")
	}

	service, err := InjectService(&Config{
		IsProduction:  common.IsProduction(isProduction),
		SessionKey:    "sessions",
		SessionSecret: common.SessionSecret(secret),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to inject API: %v", err))
	}

	api := service.API

	addr, ok := os.LookupEnv("ADDR")
	if !ok {
		panic("ADDR is not set")
	}

	err = api.Start(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to start API: %v", err))
	}
}
