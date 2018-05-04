package users_test

import (
	"log"
	"os"
	"testing"
	userMock "github.com/xesina/golang-realworld/mock/users"
	"github.com/xesina/golang-realworld/users"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/xesina/golang-realworld/pkg/types"
)

var (
	userRepository *userMock.UserRepository
)

func init() {
}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
func setup() {
	log.Println("Running Setup ...")
	// make sure that mocks satisfies the interfaces
	var _ users.UserRepository = (*userMock.UserRepository)(nil)
	userRepository = &userMock.UserRepository{}
}
func tearDown() {
}

func TestFindUser(t *testing.T) {
	expected := &users.User{
		ID:       1,
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}
	userRepository.On("Find", expected.ID).Return(expected, nil)
	i := users.NewUserInteractor(userRepository)
	actual, err := i.Find(expected.ID)
	require.NoError(t, err)
	assert.Equal(t, actual, expected)
	userRepository.AssertExpectations(t)
}

func TestRegisterUser(t *testing.T) {
	actual := &users.User{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}
	expected := &users.User{
		ID:       1,
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}
	userRepository.On("Create", actual).Return(expected, nil)
	i := users.NewUserInteractor(userRepository)
	u, err := i.Register(actual.Username, actual.Email, actual.Password)
	require.NoError(t, err)
	assert.Equal(t, expected, u)
	userRepository.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	expected := &users.User{
		ID:       1,
		Username: "test",
		Email:    "updated@email.com",
		Password: "test",
		Bio:      types.NewNullString("updated bio"),
		Image:    types.NewNullString("updated image"),
	}
	userRepository.On("Update", expected).Return(nil)
	i := users.NewUserInteractor(userRepository)
	err := i.Update(expected)
	require.NoError(t, err)
	userRepository.AssertExpectations(t)
}
