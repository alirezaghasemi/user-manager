package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/alirezaghasemi/user-manager/internal/entities"
	"github.com/alirezaghasemi/user-manager/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint64) (entities.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, user entities.User) (entities.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user entities.User) (entities.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestUserUsecase_FindByID(t *testing.T) {
	// make instance of validator
	validate := validator.New()

	mockRepo := &MockUserRepository{}

	userUsecase := NewUserUsecase(mockRepo, validate)

	// definition variables for test
	ctx := context.Background()
	userID := uint64(1)
	expectedUser := entities.User{
		ID:     userID,
		Name:   "Test User",
		Family: "Test Family",
		Email:  "test@gmail.com",
		Age:    31,
	}

	// setup expectations
	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, userID).Return(expectedUser, nil).Once()

		user, err := userUsecase.FindByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, userID).Return(entities.User{}, repository.ErrMsgUserNotFound).Once()

		user, err := userUsecase.FindByID(ctx, userID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrMsgUserNotFound))
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		genericError := errors.New("database error")

		mockRepo.On("FindByID", ctx, userID).Return(entities.User{}, genericError).Once()

		user, err := userUsecase.FindByID(ctx, userID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrMsgInternalServerError))
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_FindAll(t *testing.T) {
	validate := validator.New()

	mockRepo := &MockUserRepository{}

	userUsecase := NewUserUsecase(mockRepo, validate)

	ctx := context.Background()

	expectedUsers := []entities.User{
		{
			ID:     1,
			Name:   "Test User 1",
			Family: "Test Family 1",
			Email:  "test1@gmail.com",
			Age:    31,
		},
		{
			ID:     2,
			Name:   "Test User 2",
			Family: "Test Family 2",
			Email:  "test2@gmail.com",
			Age:    3,
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindAll", ctx).Return(expectedUsers, nil).Once()

		users, err := userUsecase.FindAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		genericError := errors.New("database error")

		mockRepo.On("FindAll", ctx).Return([]entities.User{}, genericError).Once()

		users, err := userUsecase.FindAll(ctx)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrMsgInternalServerError))
		assert.Equal(t, []entities.User{}, users)
		mockRepo.AssertExpectations(t)
	})
}
