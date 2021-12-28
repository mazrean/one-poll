package gorm2

import "github.com/google/uuid"

var (
	tables = []interface{}{
		&UserTable{},
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
