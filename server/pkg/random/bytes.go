package random

import (
	crand "crypto/rand"
	"fmt"
)

func Secure(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := crand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}

	return b, nil
}
