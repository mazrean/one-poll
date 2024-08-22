package main

import (
	"fmt"
	"os"

	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/pkg/common"
)

func main() {
	env := os.Getenv("APP_ENV")
	isProduction := env == "production"

	secret, ok := os.LookupEnv("SESSION_SECRET")
	if !ok {
		panic("SESSION_SECRET is not set")
	}

	strRelyingPartyID, ok := os.LookupEnv("RP_ID")
	if !ok {
		strRelyingPartyID = "localhost"
	}
	relyingPartyID := values.WebAuthnRelyingPartyID(strRelyingPartyID)
	if err := relyingPartyID.Validate(); err != nil {
		panic(fmt.Sprintf("invalid RP_ID: %v", err))
	}

	strRelyingPartyOrigin, ok := os.LookupEnv("RP_ORIGIN")
	if !ok {
		strRelyingPartyOrigin = "http://localhost:8080"
	}
	relyingPartyOrigin := values.NewWebAuthnOrigin(strRelyingPartyOrigin)
	if err := relyingPartyOrigin.Validate(); err != nil {
		panic(fmt.Sprintf("invalid RP_ORIGIN: %v", err))
	}

	strRelyingPartyName, ok := os.LookupEnv("RP_NAME")
	if !ok {
		strRelyingPartyName = "One Poll"
	}
	relyingPartyName := values.NewWebAuthnRelyingPartyDisplayName(strRelyingPartyName)

	service, err := InjectService(&Config{
		IsProduction:  common.IsProduction(isProduction),
		SessionKey:    "sessions",
		SessionSecret: common.SessionSecret(secret),
		RelyingParty: domain.NewWebAuthnRelyingParty(
			relyingPartyID,
			relyingPartyOrigin,
			relyingPartyName,
		),
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
