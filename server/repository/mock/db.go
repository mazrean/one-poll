package mock

import (
	context "context"
	sql "database/sql"
	"fmt"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// DB is a mock of DB interface.
type DB struct {
	ctrl     *gomock.Controller
	recorder *DBMockRecorder
}

// DBMockRecorder is the mock recorder for MockDB.
type DBMockRecorder struct {
	mock *DB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *DB {
	mock := &DB{ctrl: ctrl}
	mock.recorder = &DBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *DB) EXPECT() *DBMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *DB) Get() (*sql.DB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*sql.DB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *DBMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*DB)(nil).Get))
}

// Transaction mocks base method.
func (m *DB) Transaction(ctx context.Context, txOpt *sql.TxOptions, fn func(context.Context) error) error {
	err := fn(ctx)
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}
