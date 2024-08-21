package repository

import (
	"errors"
	"fmt"
)

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrNoRecordDeleted       = errors.New("no record deleted")
	ErrNoRecordUpdated       = errors.New("no record updated")
	ErrInvalidParameterValue = func(paramName, constraint string) error {
		return fmt.Errorf("Invalid Parameter Value : parameter '%s' must %s", paramName, constraint)
	}
	ErrDuplicateRecord = errors.New("duplicate record")
)
