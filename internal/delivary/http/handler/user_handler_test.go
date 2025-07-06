package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alirezaghasemi/user-manager/internal/delivary/http/dto/response"
	"github.com/alirezaghasemi/user-manager/internal/delivary/http/handler"
	"github.com/alirezaghasemi/user-manager/internal/entities"
	httpresponse "github.com/alirezaghasemi/user-manager/internal/pkg/httpResponse"
	"github.com/alirezaghasemi/user-manager/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Create(ctx context.Context, user entities.User) (entities.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserUsecase) FindByID(ctx context.Context, id uint64) (entities.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserUsecase) Update(ctx context.Context, user entities.User) (entities.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserUsecase) FindAll(ctx context.Context) ([]entities.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *MockUserUsecase) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_FindByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validate := validator.New()

	mockUsecase := &MockUserUsecase{}

	userHandler := handler.NewUserHandler(mockUsecase, validate)

	// ctx := context.Background()
	userID := uint64(1)

	expectedUser := entities.User{
		ID:     userID,
		Name:   "Test User",
		Family: "Test Family",
		Email:  "test@gmail.com",
		Age:    31,
	}

	setupGinContext := func(t *testing.T, path string, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, path, nil)
		c.Params = params
		return c, w
	}

	t.Run("Success", func(t *testing.T) {
		mockUsecase.On("FindByID", mock.Anything, userID).Return(expectedUser, nil).Once()

		// create http request
		// req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", userID), nil)
		// w := httptest.NewRecorder()
		// c, _ := gin.CreateTestContext(w)
		// c.Request = req
		// c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", userID)}}

		c, w := setupGinContext(t, fmt.Sprintf("/users/%d", userID), gin.Params{{Key: "id", Value: fmt.Sprintf("%d", userID)}})

		// run handler
		userHandler.FindByID(c)

		// process result
		assert.Equal(t, http.StatusOK, w.Code)

		var res httpresponse.APIResponse
		err := json.NewDecoder(w.Body).Decode(&res)

		assert.NoError(t, err)
		assert.True(t, res.Success, "response should be successful")
		assert.Equal(t, handler.SuccessMsgFoundUserById, res.Message)

		// تبدیل res.Data به response.FindUserByIDResponse
		dataBytes, err := json.Marshal(res.Data)
		assert.NoError(t, err)
		var actualResponse response.FindUserByIDResponse
		err = json.Unmarshal(dataBytes, &actualResponse)
		assert.NoError(t, err)

		expectedResponse := response.FindUserByIDResponse{
			ID:     expectedUser.ID,
			Name:   expectedUser.Name,
			Family: expectedUser.Family,
			Email:  expectedUser.Email,
			Age:    expectedUser.Age,
		}

		assert.Equal(t, expectedResponse, actualResponse)
		assert.Nil(t, res.Errors, "errors should be nil")
		mockUsecase.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		c, w := setupGinContext(t, "/users/invalid", gin.Params{{Key: "id", Value: "invalid"}})

		userHandler.FindByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var res httpresponse.APIResponse
		err := json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.False(t, res.Success, "response should not be successful")
		assert.Equal(t, handler.ErrMsgInvalidId.Error(), res.Message)
		assert.NotNil(t, res.Errors, "errors should not be nil")
		// assert.Contains(t, fmt.Sprintf("%v", res.Errors), "invalid id")
		assert.Nil(t, res.Data, "data should be nil")
		mockUsecase.AssertNotCalled(t, "FindByID", mock.Anything, "invalid")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// تنظیم رفتار Mock برای خطای UserNotFound
		mockUsecase.On("FindByID", mock.Anything, userID).Return(entities.User{}, usecase.ErrMsgUserNotFound).Once()

		// ایجاد درخواست HTTP
		c, w := setupGinContext(t, fmt.Sprintf("/users/%d", userID), gin.Params{{Key: "id", Value: fmt.Sprintf("%d", userID)}})

		// اجرای Handler
		userHandler.FindByID(c)

		// بررسی نتایج
		assert.Equal(t, http.StatusNotFound, w.Code)
		var res httpresponse.APIResponse
		err := json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.False(t, res.Success, "response should not be successful")
		assert.Equal(t, handler.ErrMsgUserNotFound.Error(), res.Message)
		assert.NotNil(t, res.Errors, "errors should not be nil")
		// assert.Contains(t, fmt.Sprintf("%v", res.Errors), "user not found")
		assert.Nil(t, res.Data, "data should be nil")
		mockUsecase.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		// تعریف خطای عمومی
		genericError := errors.New("database error")

		// تنظیم رفتار Mock برای خطای عمومی
		mockUsecase.On("FindByID", mock.Anything, userID).Return(entities.User{}, genericError).Once()

		// ایجاد درخواست HTTP
		c, w := setupGinContext(t, fmt.Sprintf("/users/%d", userID), gin.Params{{Key: "id", Value: fmt.Sprintf("%d", userID)}})

		// اجرای Handler
		userHandler.FindByID(c)

		// بررسی نتایج
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var res httpresponse.APIResponse
		err := json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.False(t, res.Success, "response should not be successful")
		assert.Equal(t, handler.ErrMsgInternalServerError.Error(), res.Message)
		assert.NotNil(t, res.Errors, "errors should not be nil")
		// assert.Contains(t, fmt.Sprintf("%v", res.Errors), "database error")
		assert.Nil(t, res.Data, "data should be nil")
		mockUsecase.AssertExpectations(t)
	})
}
