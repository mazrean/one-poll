package gorm2

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/one-poll/domain"
	"github.com/mazrean/one-poll/domain/values"
	"github.com/mazrean/one-poll/repository"
	"gorm.io/gorm"
)

type User struct {
	db *DB
}

func NewUser(db *DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) CreateUser(ctx context.Context, user *domain.User) error {
	db, err := u.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	userTable := UserTable{
		ID:       uuid.UUID(user.GetID()),
		Name:     string(user.GetName()),
		Password: string(user.GetHashedPassword()),
	}
	err = db.Create(&userTable).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (u *User) DeleteUser(ctx context.Context, id values.UserID) error {
	db, err := u.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	result := db.
		Where("id = ?", uuid.UUID(id)).
		Delete(UserTable{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected == 0 {
		return repository.ErrNoRecordDeleted
	}

	return nil
}

func (u *User) GetUser(ctx context.Context, userID values.UserID, lockType repository.LockType) (*domain.User, error) {
	db, err := u.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = u.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var userTable UserTable
	err = db.
		Where("id = ?", uuid.UUID(userID)).
		Take(&userTable).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := domain.NewUser(
		values.NewUserIDFromUUID(userTable.ID),
		values.NewUserName(userTable.Name),
		values.NewUserHashedPassword([]byte(userTable.Password)),
	)

	return user, nil
}

func (u *User) GetUserByName(ctx context.Context, name values.UserName) (*domain.User, error) {
	db, err := u.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var userTable UserTable
	err = db.
		Where("name = ?", string(name)).
		Take(&userTable).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := domain.NewUser(
		values.NewUserIDFromUUID(userTable.ID),
		values.NewUserName(userTable.Name),
		values.NewUserHashedPassword([]byte(userTable.Password)),
	)

	return user, nil
}
