package gorm2

type Tag struct {
	db *DB
}

func NewTag(db *DB) *Tag {
	return &Tag{
		db: db,
	}
}
