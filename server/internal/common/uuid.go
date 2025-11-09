package common

import "github.com/google/uuid"

// UUIDGenerator структура работающая с генерацией uuid
type UUIDGenerator struct{}

func (UUIDGenerator) New() uuid.UUID {
	return uuid.New()
}
