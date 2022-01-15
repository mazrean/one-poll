package gorm2

type Choice struct {
	db *DB
}

func NewChoice(db *DB) *Choice {
	return &Choice{
		db: db,
	}
}
