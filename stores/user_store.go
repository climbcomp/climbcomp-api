package stores

import (
	"errors"
	"fmt"
	"sync"

	"github.com/climbcomp/climbcomp-api/models"

	"github.com/gofrs/uuid"
	"github.com/jonboulle/clockwork"
)

// UserStore interface for persisting users
type UserStore interface {
	All() []models.User
	Find(id uuid.UUID) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(user models.User) (models.User, error)
	Reset()
}

// In Memory implementation of the UserStore interface
type UserMemoryStore struct {
	Clock  clockwork.Clock
	lookup map[uuid.UUID]models.User
	mux    sync.RWMutex
}

// Constructs a new UserMemoryStore instance
func NewUserMemoryStore(clock clockwork.Clock) *UserMemoryStore {
	return &UserMemoryStore{
		Clock:  clock,
		lookup: make(map[uuid.UUID]models.User),
	}
}

func (s *UserMemoryStore) All() []models.User {
	s.mux.RLock()
	defer s.mux.RUnlock()

	users := make([]models.User, 0, len(s.lookup))

	for _, value := range s.lookup {
		users = append(users, value)
	}

	return users
}

func (s *UserMemoryStore) Find(id uuid.UUID) (models.User, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	user, found := s.lookup[id]
	if !found {
		return user, fmt.Errorf("Not found: %d", id)
	}

	return user, nil
}

func (s *UserMemoryStore) Create(u models.User) (models.User, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	user := u.Clone()
	user.EnsureID()

	for _, existing := range s.lookup {
		if user.ID == existing.ID {
			return user, errors.New("ID must be unique")
		}
		if user.Name == existing.Name {
			return user, errors.New("Name must be unique")
		}
		if user.Email == existing.Email {
			return user, errors.New("Email must be unique")
		}
	}

	user.CreatedTime = s.Clock.Now()
	user.UpdatedTime = s.Clock.Now()

	s.lookup[user.ID] = user

	return user, nil
}

func (s *UserMemoryStore) Update(u models.User) (models.User, error) {
	// Ensure they're in the store first
	user, err := s.Find(u.ID)
	if err != nil {
		return user, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	// TODO: update and replace user in store

	return user, nil
}
func (s *UserMemoryStore) Delete(u models.User) (models.User, error) {
	// Ensure they're in the store first
	user, err := s.Find(u.ID)
	if err != nil {
		return user, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.lookup, u.ID)

	return user, nil
}

// Helper for resetting state in between tests
func (s *UserMemoryStore) Reset() {
	s.lookup = make(map[uuid.UUID]models.User)
}