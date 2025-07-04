package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alirezaghasemi/user-manager/internal/entities"
	"github.com/alirezaghasemi/user-manager/internal/repository"
	"github.com/alirezaghasemi/user-manager/internal/usecase"
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
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserUsecase_FindByID(t *testing.T) {
	// make instance of validator
	validate := validator.New()

	mockRepo := &MockUserRepository{}

	userUsecase := usecase.NewUserUsecase(mockRepo, validate)

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
		assert.True(t, errors.Is(err, usecase.ErrMsgUserNotFound))
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, userID).Return(entities.User{}, repository.ErrMsgInternalServerError).Once()

		user, err := userUsecase.FindByID(ctx, userID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgInternalServerError))
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_FindAll(t *testing.T) {
	validate := validator.New()

	mockRepo := &MockUserRepository{}

	userUsecase := usecase.NewUserUsecase(mockRepo, validate)

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
		mockRepo.On("FindAll", ctx).Return([]entities.User{}, repository.ErrMsgInternalServerError).Once()

		users, err := userUsecase.FindAll(ctx)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgInternalServerError))
		assert.Equal(t, []entities.User{}, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Create(t *testing.T) {
	validate := validator.New()

	mockRepo := &MockUserRepository{}

	userUsecase := usecase.NewUserUsecase(mockRepo, validate)

	ctx := context.Background()
	createUser := entities.User{
		ID:     1,
		Name:   "Test User 1",
		Family: "Test Family 1",
		Email:  "test1@gmail.com",
		Age:    31,
	}

	createdUser := entities.User{
		ID:     1,
		Name:   "Test User 1",
		Family: "Test Family 1",
		Email:  "test1@gmail.com",
		Age:    31,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Save", ctx, createUser).Return(createdUser, nil).Once()

		user, err := userUsecase.Create(ctx, createUser)

		assert.NoError(t, err)
		assert.Equal(t, createdUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DuplicateUser", func(t *testing.T) {
		mockRepo.On("Save", ctx, createUser).Return(entities.User{}, repository.ErrMsgDuplicateUser).Once()

		_, err := userUsecase.Create(ctx, createUser)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgDuplicateUser), "error should wrap usecase.ErrMsgDuplicateUser")
		assert.True(t, errors.Is(err, repository.ErrMsgDuplicateUser), "error should wrap repository.ErrMsgDuplicateUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("FailedToSaveUser", func(t *testing.T) {
		mockRepo.On("Save", ctx, createUser).Return(entities.User{}, repository.ErrMsgFailedToSaveUser).Once()

		_, err := userUsecase.Create(ctx, createUser)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgFailedToSaveUser), "error should wrap usecase.ErrMsgFailedToSaveUser")
		assert.True(t, errors.Is(err, repository.ErrMsgFailedToSaveUser), "error should wrap repository.ErrMsgFailedToSaveUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRepo.On("Save", ctx, createUser).Return(entities.User{}, repository.ErrMsgInternalServerError).Once()

		_, err := userUsecase.Create(ctx, createUser)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgInternalServerError), "error should wrap usecase.ErrMsgFailedToSaveUser")
		assert.True(t, errors.Is(err, repository.ErrMsgInternalServerError), "error should wrap repository.ErrMsgFailedToSaveUser")
		mockRepo.AssertExpectations(t)
	})

}

func TestUserUsecase_Update(t *testing.T) {
	validate := validator.New()

	mockRepo := &MockUserRepository{}

	userUsecase := usecase.NewUserUsecase(mockRepo, validate)

	ctx := context.Background()

	updateUser := entities.User{
		ID:     1,
		Name:   "Test User 1",
		Family: "Test Family 1",
		Email:  "test1@gmail.com",
		Age:    31,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Update", ctx, updateUser).Return(updateUser, nil).Once()

		user, err := userUsecase.Update(ctx, updateUser)

		assert.NoError(t, err)
		assert.Equal(t, updateUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DuplicateUser", func(t *testing.T) {
		mockRepo.On("Update", ctx, updateUser).Return(entities.User{}, repository.ErrMsgDuplicateUser).Once()

		user, err := userUsecase.Update(ctx, updateUser)

		assert.Error(t, err)
		assert.ErrorIs(t, err, usecase.ErrMsgDuplicateUser, "error should wrap usecase.ErrMsgDuplicateUser")
		// assert.True(t, errors.Is(err, usecase.ErrMsgDuplicateUser), "error should wrap usecase.ErrMsgDuplicateUser")
		assert.True(t, errors.Is(err, repository.ErrMsgDuplicateUser), "error should wrap repository.ErrMsgDuplicateUser")
		assert.Contains(t, err.Error(), "duplicate user") // check error text
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FailedToUpdateUser", func(t *testing.T) {
		mockRepo.On("Update", ctx, updateUser).Return(entities.User{}, repository.ErrMsgFailedToUpdateUser).Once()

		user, err := userUsecase.Update(ctx, updateUser)

		assert.Error(t, err)
		assert.ErrorIs(t, err, usecase.ErrMsgFailedToUpdateUser, "error should wrap usecase.ErrMsgFailedToUpdateUser")
		// assert.True(t, errors.Is(err, usecase.ErrMsgFailedToUpdateUser), "error should wrap usecase.ErrMsgFailedToUpdateUser")
		assert.True(t, errors.Is(err, repository.ErrMsgFailedToUpdateUser), "error should wrap repository.ErrMsgFailedToUpdateUser")
		assert.Contains(t, err.Error(), "failed to update user")
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.On("Update", ctx, updateUser).Return(entities.User{}, repository.ErrMsgUserNotFound).Once()

		user, err := userUsecase.Update(ctx, updateUser)

		assert.Error(t, err)
		assert.ErrorIs(t, err, usecase.ErrMsgUserNotFound, "error should be usecase.ErrMsgUserNotFound")
		// assert.True(t, errors.Is(err, usecase.ErrMsgUserNotFound), "error should wrap usecase.ErrMsgUserNotFound")
		assert.True(t, errors.Is(err, repository.ErrMsgUserNotFound), "error should wrap repository.ErrMsgUserNotFound")
		assert.Contains(t, err.Error(), "user not found")
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		genericError := errors.New("database error")

		mockRepo.On("Update", ctx, updateUser).Return(entities.User{}, genericError).Once()

		user, err := userUsecase.Update(ctx, updateUser)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgInternalServerError), "error should wrap usecase.ErrMsgInternalServerError")
		assert.Contains(t, err.Error(), "internal server")
		assert.Equal(t, entities.User{}, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Delete(t *testing.T) {
	validate := validator.New()

	mockRepo := &MockUserRepository{}

	userUsecase := usecase.NewUserUsecase(mockRepo, validate)

	ctx := context.Background()
	userID := uint64(1)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, userID).Return(nil).Once()

		err := userUsecase.Delete(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.On("Delete", ctx, userID).Return(repository.ErrMsgUserNotFound).Once()

		err := userUsecase.Delete(ctx, userID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgUserNotFound), "error should wrap usecase.ErrMsgUserNotFound")
		assert.True(t, errors.Is(err, repository.ErrMsgUserNotFound), "error should wrap repository.ErrMsgUserNotFound")
		mockRepo.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRepo.On("Delete", ctx, userID).Return(repository.ErrMsgInternalServerError).Once()

		err := userUsecase.Delete(ctx, userID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, usecase.ErrMsgInternalServerError), "error should wrap usecase.ErrMsgInternalServerError")
		assert.True(t, errors.Is(err, repository.ErrMsgInternalServerError), "error should wrap repository.ErrMsgInternalServerError")
		mockRepo.AssertExpectations(t)
	})
}
