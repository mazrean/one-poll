package values

import "github.com/google/uuid"

type (
	ResponseID uuid.UUID
)

func NewResponseID() ResponseID {
	return ResponseID(uuid.New())
}

func NewResponseIDFromUUID(id uuid.UUID) ResponseID {
	return ResponseID(id)
}
