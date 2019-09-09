package stores_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/climbcomp/climbcomp-api/models"
	"github.com/climbcomp/climbcomp-api/stores"

	"github.com/gofrs/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Run the test suite for the in-memory implementation
func TestUserMemoryStore(t *testing.T) {
	clock := clockwork.NewFakeClock()
	suite.Run(t, &UserStoreTestSuite{
		Clock: clock,
		Store: stores.NewUserMemoryStore(clock),
	})
}

// Run the test suite for the postgres implementation
// TODO

// Test suite for the full UserStore interface
type UserStoreTestSuite struct {
	suite.Suite
	Clock clockwork.FakeClock
	Store stores.UserStore
}

func (s *UserStoreTestSuite) SetupTest() {
}

func (s *UserStoreTestSuite) TearDownTest() {
	s.Store.Reset()
}

// Helpers

func (s *UserStoreTestSuite) CreateUser(name string) models.User {
	return s.CreateUserWithID(models.NewUUID(), name)
}

func (s *UserStoreTestSuite) CreateUserWithID(id uuid.UUID, name string) models.User {
	user, err := s.Store.Create(models.User{
		ID:    id,
		Email: fmt.Sprintf("%v@example.com", name),
		Name:  name,
	})
	assert.NoError(s.T(), err)
	return user
}

// Tests

func (s *UserStoreTestSuite) TestAllReturnsAllEntities() {
	t := s.T()

	assert.Equal(t, 0, len(s.Store.All()))

	_ = s.CreateUser("user1")
	assert.Equal(t, 1, len(s.Store.All()))

	_ = s.CreateUser("user2")
	assert.Equal(t, 2, len(s.Store.All()))
}

func (s *UserStoreTestSuite) TestFind() {
	t := s.T()
	uuid := models.NewUUID()

	_, err := s.Store.Find(uuid)
	assert.EqualError(t, err, "ID not found")

	_ = s.CreateUserWithID(uuid, "user1")

	_, err = s.Store.Find(uuid)
	assert.NoError(t, err)
}

func (s *UserStoreTestSuite) TestCreate() {
	t := s.T()
	now := s.Clock.Now()
	uuid := models.NewUUID()

	found, _ := s.Store.Exists(uuid)
	assert.False(t, found)

	user, err := s.Store.Create(models.User{
		ID:    uuid,
		Email: "user1@example.com",
		Name:  "user1",
	})
	assert.NoError(t, err)
	assert.Equal(t, uuid, user.ID)
	assert.Equal(t, now, user.CreatedTime)
	assert.Equal(t, now, user.UpdatedTime)

	found, _ = s.Store.Exists(uuid)
	assert.True(t, found)
}

func (s *UserStoreTestSuite) TestCreateWithNoID() {
	t := s.T()
	now := s.Clock.Now()

	user, err := s.Store.Create(models.User{
		Email: "user1@example.com",
		Name:  "user1",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, now, user.CreatedTime)
	assert.Equal(t, now, user.UpdatedTime)

	found, _ := s.Store.Exists(user.ID)
	assert.True(t, found)
}

func (s *UserStoreTestSuite) TestCreateWithConstraintViolations() {
	t := s.T()
	uuid1 := models.NewUUID()
	uuid2 := models.NewUUID()

	_, err := s.Store.Create(models.User{
		ID:    uuid1,
		Email: "user1@example.com",
		Name:  "user1",
	})
	assert.NoError(t, err)

	_, err = s.Store.Create(models.User{
		ID:    uuid1,
		Email: "user2@example.com",
		Name:  "user2",
	})
	assert.EqualError(t, err, "ID must be unique")

	_, err = s.Store.Create(models.User{
		ID:    uuid2,
		Email: "user1@example.com",
		Name:  "user2",
	})
	assert.EqualError(t, err, "Email must be unique")

	_, err = s.Store.Create(models.User{
		ID:    uuid2,
		Email: "user2@example.com",
		Name:  "user1",
	})
	assert.EqualError(t, err, "Name must be unique")
}

func (s *UserStoreTestSuite) TestUpdate() {
	t := s.T()
	createdTime := s.Clock.Now()
	user := s.CreateUser("user1")

	user.Name = "user1.1"
	user.Email = "user1.1@example.com"

	s.Clock.Advance(1 * time.Hour)
	updatedTime := s.Clock.Now()
	_, err := s.Store.Update(user)
	assert.NoError(t, err)

	persisted, err := s.Store.Find(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "user1.1", persisted.Name)
	assert.Equal(t, "user1.1@example.com", persisted.Email)
	assert.Equal(t, createdTime, persisted.CreatedTime)
	assert.Equal(t, updatedTime, persisted.UpdatedTime)
}

func (s *UserStoreTestSuite) TestDelete() {
	t := s.T()
	uuid := models.NewUUID()

	_ = s.CreateUserWithID(uuid, "user1")
	found, _ := s.Store.Exists(uuid)
	assert.True(t, found)

	deleted, err := s.Store.Delete(uuid)
	assert.True(t, deleted)
	assert.NoError(t, err)

	found, _ = s.Store.Exists(uuid)
	assert.False(t, found)

	deleted, err = s.Store.Delete(uuid)
	assert.False(t, deleted)
	assert.EqualError(t, err, "ID not found")
}
