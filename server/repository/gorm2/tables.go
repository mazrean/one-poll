package gorm2

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	tables = []interface{}{
		&UserTable{},
		&PollTable{},
		&PollTypeTable{},
		&TagTable{},
		&ChoiceTable{},
		&Response{},
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

type PollTable struct {
	ID        uuid.UUID      `gorm:"type:char(36);not null;primaryKey;size:36"`
	OwnerID   uuid.UUID      `gorm:"type:char(36);not null;size:36"`
	Title     string         `gorm:"type:varchar(50);not null;size:50"`
	TypeID    int            `gorm:"type:int(11);not null"`
	Deadline  time.Time      `gorm:"type:DATETIME NULL;default:NULL"`
	CreatedAt time.Time      `gorm:"type:datetime;not null"`
	DeletedAt gorm.DeletedAt `gorm:"type:DATETIME NULL;default:NULL"`
	Owner     UserTable      `gorm:"foreignKey:OwnerID"`
	Choices   []ChoiceTable  `gorm:"foreignKey:PollID"`
	Tags      []TagTable     `gorm:"many2many:poll_tag_relations"`
}

func (*PollTable) TableName() string {
	return "polls"
}

type PollTypeTable struct {
	ID     int    `gorm:"type:int(11);not null;primaryKey;autoIncrement"`
	Name   string `gorm:"type:varchar(50);not null;size:50;unique"`
	Active bool   `gorm:"type:bool;not null;default:true"`
}

func (*PollTypeTable) TableName() string {
	return "poll_types"
}

type TagTable struct {
	ID   int    `gorm:"type:int(11);not null;primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(50);not null;size:50;unique"`
}

func (*TagTable) TableName() string {
	return "tags"
}

type ChoiceTable struct {
	ID     uuid.UUID `gorm:"type:char(36);not null;primaryKey;size:36"`
	PollID uuid.UUID `gorm:"type:char(36);not null;size:36"`
	Name   string    `gorm:"type:varchar(50);not null;size:50"`
}

func (*ChoiceTable) TableName() string {
	return "choices"
}

type Response struct {
	ID           uuid.UUID     `gorm:"type:char(36);not null;primaryKey;size:36"`
	PollID       uuid.UUID     `gorm:"type:char(36);not null;size:36"`
	RespondentID uuid.UUID     `gorm:"type:char(36);not null;size:36"`
	CreatedAt    time.Time     `gorm:"type:datetime;not null"`
	Poll         PollTable     `gorm:"foreignKey:PollID"`
	Respondent   UserTable     `gorm:"foreignKey:RespondentID"`
	Choices      []ChoiceTable `gorm:"many2many:response_choice_relations"`
	Comment      Comment       `gorm:"foreignKey:ResponseID"`
}

func (*Response) TableName() string {
	return "responses"
}
