package v1

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
)

type Checker struct{}

func NewChecker() *Checker {
	return &Checker{}
}

func (m *Checker) check(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	checkerMap := map[string]openapi3filter.AuthenticationFunc{}

	checker, ok := checkerMap[input.SecuritySchemeName]
	if !ok {
		return nil
		// checkerが出揃うまで一時的にエラーにしない
		// return fmt.Errorf("unknown security scheme: %q", input.SecuritySchemeName)
	}

	return checker(ctx, input)
}
