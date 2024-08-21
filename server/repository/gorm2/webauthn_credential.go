package gorm2

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
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

func (wac *WebAuthnCredential) StoreCredential(ctx context.Context, userID values.UserID, credential *domain.WebAuthnCredential) error {
	db, err := wac.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	algorismID, ok := wac.algorismMap[credential.Algorithm()]
	if !ok {
		return fmt.Errorf("unknown algorism: %d", credential.Algorithm())
	}

	transports := make([]WebAuthnCredentialTransportTable, 0, len(credential.Transports()))
	for _, transport := range credential.Transports() {
		transportType, ok := wac.transportMap[transport]
		if !ok {
			return fmt.Errorf("unknown transport: %s", transport)
		}

		transports = append(transports, WebAuthnCredentialTransportTable{
			CredentialID: uuid.UUID(credential.ID()),
			TypeID:       transportType.ID,
		})
	}

	credentialTable := WebAuthnCredentialTable{
		ID:          uuid.UUID(credential.ID()),
		UserID:      uuid.UUID(userID),
		Name:        string(credential.Name()),
		AlgorithmID: algorismID.ID,
		Transports:  transports,
		CreatedAt:   credential.CreatedAt(),
		LastUsedAt:  credential.LastUsedAt(),
	}

	err = db.Create(&credentialTable).Error
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	return nil
}

func (wac *WebAuthnCredential) GetCredentialsByUserID(ctx context.Context, userID values.UserID) ([]*domain.WebAuthnCredential, error) {
	db, err := wac.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var credentialTables []WebAuthnCredentialTable
	err = db.
		Joins("Algorithm").
		Preload("Transports.Type").
		Where("user_id = ?", uuid.UUID(userID)).
		Find(&credentialTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	credentials := make([]*domain.WebAuthnCredential, 0, len(credentialTables))
CREDENTIAL_LOOP:
	for _, credentialTable := range credentialTables {
		var algorithm values.WebAuthnCredentialAlgorithm
		switch credentialTable.Algorithm.Name {
		case algorismES256:
			algorithm = values.WebAuthnCredentialAlgorithmES256
		default:
			log.Printf("error: unknown algorism: %s", credentialTable.Algorithm.Name)
			continue
		}

		transports := make([]values.WebAuthnCredentialTransport, 0, len(credentialTable.Transports))
		for _, transport := range credentialTable.Transports {
			var transportType values.WebAuthnCredentialTransport
			switch transport.Type.Name {
			case transportUSB:
				transportType = values.WebAuthnCredentialTransportUSB
			case transportNFC:
				transportType = values.WebAuthnCredentialTransportNFC
			case transportBLE:
				transportType = values.WebAuthnCredentialTransportBLE
			case transportSmartCard:
				transportType = values.WebAuthnCredentialTransportSmartCard
			case transportHybrid:
				transportType = values.WebAuthnCredentialTransportHybrid
			case transportInternal:
				transportType = values.WebAuthnCredentialTransportInternal
			default:
				log.Printf("error: unknown transport: %s", transport.Type.Name)
				continue CREDENTIAL_LOOP
			}

			transports = append(transports, transportType)
		}

		credential := domain.NewWebAuthnCredential(
			values.NewWebAuthnCredentialIDFromUUID(credentialTable.ID),
			values.NewWebAuthnCredentialCredID(credentialTable.CredID),
			values.NewWebAuthnCredentialName(credentialTable.Name),
			values.NewWebAuthnCredentialPublicKey(credentialTable.PublicKey),
			algorithm,
			transports,
			credentialTable.CreatedAt,
			credentialTable.LastUsedAt,
		)

		credentials = append(credentials, credential)
	}

	return credentials, nil
}

func (wac *WebAuthnCredential) GetCredentialWithUserByCredID(ctx context.Context, credID values.WebAuthnCredentialCredID, lockType repository.LockType) (*domain.WebAuthnCredential, *domain.User, error) {
	db, err := wac.db.getDB(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = wac.db.setLock(db, lockType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var credentialTable WebAuthnCredentialTable
	err = db.
		Joins("Algorithm").
		Joins("User").
		Preload("Transports.Type").
		Where("cred_id = ?", credID).
		First(&credentialTable).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, repository.ErrRecordNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get credential: %w", err)
	}

	var algorithm values.WebAuthnCredentialAlgorithm
	switch credentialTable.Algorithm.Name {
	case algorismES256:
		algorithm = values.WebAuthnCredentialAlgorithmES256
	default:
		return nil, nil, fmt.Errorf("unknown algorism: %s", credentialTable.Algorithm.Name)
	}

	transports := make([]values.WebAuthnCredentialTransport, 0, len(credentialTable.Transports))
	for _, transport := range credentialTable.Transports {
		var transportType values.WebAuthnCredentialTransport
		switch transport.Type.Name {
		case transportUSB:
			transportType = values.WebAuthnCredentialTransportUSB
		case transportNFC:
			transportType = values.WebAuthnCredentialTransportNFC
		case transportBLE:
			transportType = values.WebAuthnCredentialTransportBLE
		case transportSmartCard:
			transportType = values.WebAuthnCredentialTransportSmartCard
		case transportHybrid:
			transportType = values.WebAuthnCredentialTransportHybrid
		case transportInternal:
			transportType = values.WebAuthnCredentialTransportInternal
		default:
			return nil, nil, fmt.Errorf("unknown transport: %s", transport.Type.Name)
		}

		transports = append(transports, transportType)
	}

	credential := domain.NewWebAuthnCredential(
		values.NewWebAuthnCredentialIDFromUUID(credentialTable.ID),
		values.NewWebAuthnCredentialCredID(credentialTable.CredID),
		values.NewWebAuthnCredentialName(credentialTable.Name),
		values.NewWebAuthnCredentialPublicKey(credentialTable.PublicKey),
		algorithm,
		transports,
		credentialTable.CreatedAt,
		credentialTable.LastUsedAt,
	)

	user := domain.NewUser(
		values.NewUserIDFromUUID(credentialTable.UserID),
		values.NewUserName(credentialTable.User.Name),
		values.NewUserHashedPassword([]byte(credentialTable.User.Password)),
	)

	return credential, user, nil
}

func (wac *WebAuthnCredential) UpdateLastUsedAt(ctx context.Context, credential *domain.WebAuthnCredential) error {
	db, err := wac.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	result := db.
		Model(&WebAuthnCredentialTable{}).
		Where("id = ?", uuid.UUID(credential.ID())).
		Update("LastUsedAt", credential.LastUsedAt())
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update last used at: %w", err)
	}

	if result.RowsAffected == 0 {
		return repository.ErrNoRecordUpdated
	}

	return nil
}

func (wac *WebAuthnCredential) DeleteCredential(ctx context.Context, userID values.UserID, credID values.WebAuthnCredentialCredID) error {
	db, err := wac.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	result := db.
		Where("user_id = ?", uuid.UUID(userID)).
		Where("cred_id = ?", credID).
		Delete(&WebAuthnCredentialTable{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	if result.RowsAffected == 0 {
		return repository.ErrNoRecordDeleted
	}

	return nil
}
