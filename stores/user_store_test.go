package stores_test

import (
	"testing"

	"github.com/climbcomp/climbcomp-api/models"
	"github.com/climbcomp/climbcomp-api/stores"

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
	Clock clockwork.Clock
	Store stores.UserStore
}

func (s *UserStoreTestSuite) SetupTest() {
}

func (s *UserStoreTestSuite) TearDownTest() {
	s.Store.Reset()
}

func (s *UserStoreTestSuite) TestCreate() {
	t := s.T()
	now := s.Clock.Now()

	uuid := models.NewUUID()
	entity, err := s.Store.Create(models.User{
		ID:    uuid,
		Email: "user1@example.com",
		Name:  "user1",
	})
	assert.NoError(t, err)

	assert.Equal(t, uuid, entity.ID)
	assert.Equal(t, now, entity.CreatedTime)
	assert.Equal(t, now, entity.UpdatedTime)
}

func (s *UserStoreTestSuite) TestCreateWithNoID() {
	t := s.T()
	now := s.Clock.Now()

	entity, err := s.Store.Create(models.User{
		Email: "user1@example.com",
		Name:  "user1",
	})
	assert.NoError(t, err)

	assert.NotEmpty(t, entity.ID)
	assert.Equal(t, now, entity.CreatedTime)
	assert.Equal(t, now, entity.UpdatedTime)
}

func (s *UserStoreTestSuite) TestAllReturnsAllEntities() {
	t := s.T()

	// Sanity check
	assert.Equal(t, 0, len(s.Store.All()))

	_, err := s.Store.Create(models.User{
		Email: "user1@example.com",
		Name:  "user1",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(s.Store.All()))

	_, err = s.Store.Create(models.User{
		Email: "user2@example.com",
		Name:  "user2",
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(s.Store.All()))
}
