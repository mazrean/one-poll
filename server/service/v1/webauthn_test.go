package v1

import "testing"

func TestLoadAAGUID2NameMap(t *testing.T) {
	_, err := loadAAGUID2NameMap()
	if err != nil {
		t.Errorf("loadAAGUID2NameMap() failed: %v", err)
	}
}
