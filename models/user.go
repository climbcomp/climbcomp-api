package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Represents a user of the API
type User struct {
	ID          uuid.UUID
	CreatedTime time.Time
	UpdatedTime time.Time

	Email string
	Name  string
	Slug  string
}

func (s *User) Clone() User {
	return User{
		ID:          s.ID,
		CreatedTime: s.CreatedTime,
		UpdatedTime: s.UpdatedTime,
		Email:       s.Email,
		Name:        s.Name,
		Slug:        s.Slug,
	}
}

func (s *User) EnsureID() {
	if s.ID == uuid.Nil {
		s.ID = NewUUID()
	}
}

func (s *User) HasID() bool {
	return s.ID != uuid.Nil
}
