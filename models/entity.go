package models

import (
	"github.com/gofrs/uuid"
)

func NewUUID() uuid.UUID {
	// Going to assume that we can reliably create v4 UUIDs
	// (since all they are are random numbers)
	return uuid.Must(uuid.NewV4())
}
