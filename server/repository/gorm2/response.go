package gorm2

type Response struct {
	db *DB
}

func NewResponse(db *DB) *Response {
	return &Response{
		db: db,
	}
}
