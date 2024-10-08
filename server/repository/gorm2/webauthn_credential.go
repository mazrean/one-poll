package gorm2

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
	"gorm.io/gorm"
)

type WebAuthnCredential struct {
	db          *DB
	algorismMap map[values.WebAuthnCredentialAlgorithm]WebAuthnCredentialAlgorithmTable
}

func NewWebAuthnCredential(db *DB) (*WebAuthnCredential, error) {
	algorismMap, err := setupWebAuthnCredentialAlgorismTable(db.db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup algorism table: %w", err)
	}

	return &WebAuthnCredential{
		db:          db,
		algorismMap: algorismMap,
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
		if !algorism.Active {
			continue
		}

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

func (wac *WebAuthnCredential) StoreCredential(ctx context.Context, userID values.UserID, credential *domain.WebAuthnCredential) error {
	db, err := wac.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	algorismID, ok := wac.algorismMap[credential.Algorithm()]
	if !ok {
		return fmt.Errorf("unknown algorism: %d", credential.Algorithm())
	}

	base64.StdEncoding.EncodeToString(credential.CredID())

	credentialTable := WebAuthnCredentialTable{
		ID:          uuid.UUID(credential.ID()),
		UserID:      uuid.UUID(userID),
		CredID:      base64.StdEncoding.EncodeToString(credential.CredID()),
		Name:        string(credential.Name()),
		PublicKey:   credential.PublicKey(),
		AlgorithmID: algorismID.ID,
		CreatedAt:   credential.CreatedAt(),
		LastUsedAt:  credential.LastUsedAt(),
	}

	err = db.Create(&credentialTable).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return repository.ErrDuplicateRecord
		}

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
		Where("user_id = ?", uuid.UUID(userID)).
		Where("Algorithm.active = true").
		Order("created_at DESC").
		Find(&credentialTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	credentials := make([]*domain.WebAuthnCredential, 0, len(credentialTables))
	for _, credentialTable := range credentialTables {
		var algorithm values.WebAuthnCredentialAlgorithm
		switch credentialTable.Algorithm.Name {
		case algorismES256:
			algorithm = values.WebAuthnCredentialAlgorithmES256
		default:
			log.Printf("error: unknown algorism: %s", credentialTable.Algorithm.Name)
			continue
		}

		credID, err := base64.StdEncoding.DecodeString(credentialTable.CredID)
		if err != nil {
			log.Printf("error: failed to decode credID: %v", err)
			continue
		}

		credential := domain.NewWebAuthnCredential(
			values.NewWebAuthnCredentialIDFromUUID(credentialTable.ID),
			values.NewWebAuthnCredentialCredID(credID),
			values.NewWebAuthnCredentialName(credentialTable.Name),
			values.NewWebAuthnCredentialPublicKey(credentialTable.PublicKey),
			algorithm,
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
		Where("cred_id = ?", base64.StdEncoding.EncodeToString(credID)).
		Where("Algorithm.active = true").
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

	credential := domain.NewWebAuthnCredential(
		values.NewWebAuthnCredentialIDFromUUID(credentialTable.ID),
		credID,
		values.NewWebAuthnCredentialName(credentialTable.Name),
		values.NewWebAuthnCredentialPublicKey(credentialTable.PublicKey),
		algorithm,
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

func (wac *WebAuthnCredential) DeleteCredential(ctx context.Context, userID values.UserID, credentialID values.WebAuthnCredentialID) error {
	db, err := wac.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	result := db.
		Where("user_id = ?", uuid.UUID(userID)).
		Where("id = ?", uuid.UUID(credentialID)).
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
