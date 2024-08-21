package gorm2

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	tables = []interface{}{
		&UserTable{},
		&WebAuthnCredentialTable{},
		&WebAuthnCredentialAlgorithmTable{},
		&WebAuthnCredentialTransportTypeTable{},
		&WebAuthnCredentialTransportTable{},
		&PollTable{},
		&PollTypeTable{},
		&TagTable{},
		&ChoiceTable{},
		&ResponseTable{},
		&CommentTable{},
	}
)

type UserTable struct {
	ID       uuid.UUID `gorm:"type:char(36);not null;primaryKey;size:36"`
	Name     string    `gorm:"type:varchar(50);not null;size:50"`
	Password string    `gorm:"type:char(60);not null;size:60"`
}

func (*UserTable) TableName() string {
	return "users"
}

type WebAuthnCredentialTable struct {
	ID          uuid.UUID                          `gorm:"type:char(36);not null;primaryKey;size:36"`
	UserID      uuid.UUID                          `gorm:"type:char(36);not null;size:36"`
	User        UserTable                          `gorm:"foreignKey:ID;references:UserID"`
	CredID      []byte                             `gorm:"type:binary(64);not null;size:64"`
	Name        string                             `gorm:"type:varchar(50);not null;size:50"`
	PublicKey   []byte                             `gorm:"type:binary(65);not null;size:65"`
	AlgorithmID int                                `gorm:"type:tinyint;not null"`
	Algorithm   WebAuthnCredentialAlgorithmTable   `gorm:"foreignKey:AlgorithmID"`
	Transports  []WebAuthnCredentialTransportTable `gorm:"foreignKey:CredentialID"`
	CreatedAt   time.Time                          `gorm:"type:datetime;not null"`
	LastUsedAt  time.Time                          `gorm:"type:datetime;not null"`
}

func (*WebAuthnCredentialTable) TableName() string {
	return "webauthn_credentials"
}

type WebAuthnCredentialAlgorithmTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(50);not null;size:50;unique"`
	Active bool   `gorm:"type:bool;not null;default:true"`
}

func (*WebAuthnCredentialAlgorithmTable) TableName() string {
	return "webauthn_credential_algorithms"
}

type WebAuthnCredentialTransportTypeTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(50);not null;size:50;unique"`
	Active bool   `gorm:"type:bool;not null;default:true"`
}

func (*WebAuthnCredentialTransportTypeTable) TableName() string {
	return "webauthn_credential_transport_types"
}

type WebAuthnCredentialTransportTable struct {
	CredentialID uuid.UUID `gorm:"type:char(36);not null;primaryKey;size:36"`
	TypeID       int       `gorm:"type:TINYINT;not null;primaryKey"`
	Type         WebAuthnCredentialTransportTypeTable
}

func (*WebAuthnCredentialTransportTable) TableName() string {
	return "webauthn_credential_transports"
}

type PollTable struct {
	ID        uuid.UUID      `gorm:"type:char(36);not null;primaryKey;size:36"`
	OwnerID   uuid.UUID      `gorm:"type:char(36);not null;size:36"`
	Title     string         `gorm:"type:varchar(50);not null;size:50"`
	TypeID    int            `gorm:"type:tinyint;not null"`
	Deadline  sql.NullTime   `gorm:"type:DATETIME NULL;default:NULL"`
	CreatedAt time.Time      `gorm:"type:datetime;not null"`
	DeletedAt gorm.DeletedAt `gorm:"type:DATETIME NULL;default:NULL"`
	Owner     UserTable      `gorm:"foreignKey:OwnerID"`
	Choices   []ChoiceTable  `gorm:"foreignKey:PollID"`
	Tags      []TagTable     `gorm:"many2many:poll_tag_relations"`
	PollType  PollTypeTable  `gorm:"foreignKey:TypeID"`
}

func (*PollTable) TableName() string {
	return "polls"
}

type PollTypeTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(50);not null;size:50;unique"`
	Active bool   `gorm:"type:bool;not null;default:true"`
}

func (*PollTypeTable) TableName() string {
	return "poll_types"
}

type TagTable struct {
	ID   uuid.UUID `gorm:"type:char(36);not null;primaryKey;size:36"`
	Name string    `gorm:"type:varchar(50);not null;size:50;unique"`
}

func (*TagTable) TableName() string {
	return "tags"
}

type ChoiceTable struct {
	ID     uuid.UUID `gorm:"type:char(36);not null;primaryKey;size:36"`
	PollID uuid.UUID `gorm:"type:char(36);not null;size:36"`
	Name   string    `gorm:"type:varchar(50);not null;size:50"`
	Order  uint8     `gorm:"type:TINYINT;not null;default:0"`
}

func (*ChoiceTable) TableName() string {
	return "choices"
}

type ResponseTable struct {
	ID           uuid.UUID     `gorm:"type:char(36);not null;primaryKey;size:36"`
	PollID       uuid.UUID     `gorm:"type:char(36);not null;size:36"`
	RespondentID uuid.UUID     `gorm:"type:char(36);not null;size:36"`
	CreatedAt    time.Time     `gorm:"type:datetime;not null"`
	Poll         PollTable     `gorm:"foreignKey:PollID"`
	Respondent   UserTable     `gorm:"foreignKey:RespondentID"`
	Choices      []ChoiceTable `gorm:"many2many:response_choice_relations"`
	Comment      CommentTable  `gorm:"foreignKey:ResponseID"`
}

func (*ResponseTable) TableName() string {
	return "responses"
}

type CommentTable struct {
	ID         uuid.UUID `gorm:"type:char(36);not null;primaryKey;size:36"`
	ResponseID uuid.UUID `gorm:"type:char(36);not null;size:36"`
	Comment    string    `gorm:"type:text;not null;size:5000"`
}

func (*CommentTable) TableName() string {
	return "comments"
}
