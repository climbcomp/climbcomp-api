package stores

import (
	"errors"
	"sync"

	"github.com/climbcomp/climbcomp-api/models"

	"github.com/gofrs/uuid"
	"github.com/jonboulle/clockwork"
)

// UserStore interface for persisting users
type UserStore interface {
	All() []models.User
	Exists(id uuid.UUID) (bool, error)
	Find(id uuid.UUID) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id uuid.UUID) (bool, error)
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

func (s *UserMemoryStore) Exists(id uuid.UUID) (bool, error) {
	_, err := s.Find(id)
	return err == nil, err
}

func (s *UserMemoryStore) Find(id uuid.UUID) (models.User, error) {
	user, found := s.lookup[id]
	if !found {
		return user, errors.New("ID not found")
	}
	return user, nil
}

func (s *UserMemoryStore) Create(u models.User) (models.User, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	user := u.Clone()
	user.EnsureID()

	found, _ := s.Exists(user.ID)
	if found {
		return user, errors.New("ID must be unique")
	}

	err := s.enforceConstraints(user)
	if err != nil {
		return user, err
	}

	user.CreatedTime = s.Clock.Now()
	user.UpdatedTime = s.Clock.Now()

	s.lookup[user.ID] = user

	return user, nil
}

func (s *UserMemoryStore) Update(u models.User) (models.User, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	user := u.Clone()

	found, err := s.Exists(user.ID)
	if !found {
		return user, err
	}

	err = s.enforceConstraints(user)
	if err != nil {
		return user, err
	}

	user.UpdatedTime = s.Clock.Now()

	s.lookup[u.ID] = user

	return user, nil
}

func (s *UserMemoryStore) Delete(id uuid.UUID) (bool, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	found, err := s.Exists(id)
	if !found {
		return false, err
	}

	delete(s.lookup, id)

	return true, nil
}

// Helper for resetting state in between tests
func (s *UserMemoryStore) Reset() {
	s.lookup = make(map[uuid.UUID]models.User)
}

func (s *UserMemoryStore) enforceConstraints(u models.User) error {
	for _, existing := range s.lookup {
		if u.ID == existing.ID {
			continue
		}
		if u.Name == existing.Name {
			return errors.New("Name must be unique")
		}
		if u.Email == existing.Email {
			return errors.New("Email must be unique")
		}
	}
	return nil
}
