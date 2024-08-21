package gorm2

import (
	"fmt"

	"github.com/mazrean/one-poll/domain/values"
	"gorm.io/gorm"
)

type WebAuthnCredential struct {
	db           *DB
	algorismMap  map[values.WebAuthnCredentialAlgorithm]WebAuthnCredentialAlgorithmTable
	transportMap map[values.WebAuthnCredentialTransport]WebAuthnCredentialTransportTypeTable
}

func NewWebAuthnCredential(db *DB) (*WebAuthnCredential, error) {
	algorismMap, err := setupWebAuthnCredentialAlgorismTable(db.db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup algorism table: %w", err)
	}

	transportMap, err := setupWebAuthnCredentialTransportTypeTable(db.db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup transport table: %w", err)
	}

	return &WebAuthnCredential{
		db:           db,
		algorismMap:  algorismMap,
		transportMap: transportMap,
	}, nil
}

var (
	algorismES256 = "ES256"
)

func setupWebAuthnCredentialAlgorismTable(db *gorm.DB) (map[values.WebAuthnCredentialAlgorithm]WebAuthnCredentialAlgorithmTable, error) {
	algorisms := []WebAuthnCredentialAlgorithmTable{
		{Name: algorismES256, Active: true},
	}

	for i, algorism := range algorisms {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", algorism.Name).
			FirstOrCreate(&algorism).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create resource type: %w", err)
		}

		algorisms[i] = algorism
	}

	algorismMap := make(map[values.WebAuthnCredentialAlgorithm]WebAuthnCredentialAlgorithmTable, len(algorisms))
	for _, algorism := range algorisms {
		var algorismValue values.WebAuthnCredentialAlgorithm
		switch algorism.Name {
		case algorismES256:
			algorismValue = values.WebAuthnCredentialAlgorithmES256
		default:
			return nil, fmt.Errorf("unknown algorism: %s", algorism.Name)
		}

		algorismMap[algorismValue] = algorism
	}

	return algorismMap, nil
}

var (
	transportUSB       = "USB"
	transportNFC       = "NFC"
	transportBLE       = "BLE"
	transportSmartCard = "SmartCard"
	transportHybrid    = "Hybrid"
	transportInternal  = "Internal"
)

func setupWebAuthnCredentialTransportTypeTable(db *gorm.DB) (map[values.WebAuthnCredentialTransport]WebAuthnCredentialTransportTypeTable, error) {
	transports := []WebAuthnCredentialTransportTypeTable{
		{Name: transportUSB, Active: true},
		{Name: transportNFC, Active: true},
		{Name: transportBLE, Active: true},
		{Name: transportSmartCard, Active: true},
		{Name: transportHybrid, Active: true},
		{Name: transportInternal, Active: true},
	}

	for i, transport := range transports {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", transport.Name).
			FirstOrCreate(&transport).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create resource type: %w", err)
		}

		transports[i] = transport
	}

	transportTypeMap := make(map[values.WebAuthnCredentialTransport]WebAuthnCredentialTransportTypeTable, len(transports))
	for _, transport := range transports {
		var transportValue values.WebAuthnCredentialTransport
		switch transport.Name {
		case transportUSB:
			transportValue = values.WebAuthnCredentialTransportUSB
		case transportNFC:
			transportValue = values.WebAuthnCredentialTransportNFC
		case transportBLE:
			transportValue = values.WebAuthnCredentialTransportBLE
		case transportSmartCard:
			transportValue = values.WebAuthnCredentialTransportSmartCard
		case transportHybrid:
			transportValue = values.WebAuthnCredentialTransportHybrid
		case transportInternal:
			transportValue = values.WebAuthnCredentialTransportInternal
		default:
			return nil, fmt.Errorf("unknown transport: %s", transport.Name)
		}

		transportTypeMap[transportValue] = transport
	}

	return transportTypeMap, nil
}
