package handler

import (
	"errors"
	"net/http"

	"github.com/alirezaghasemi/user-manager/internal/delivary/http/dto/request"
	"github.com/alirezaghasemi/user-manager/internal/delivary/http/dto/response"
	"github.com/alirezaghasemi/user-manager/internal/entities"
	httpresponse "github.com/alirezaghasemi/user-manager/internal/pkg/httpResponse"
	"github.com/alirezaghasemi/user-manager/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Definition Error
var (
	ErrMsgDuplicateUser       = errors.New("duplicate user error")
	ErrMsgFailedToSaveUser    = errors.New("failed to save user")
	ErrMsgInternalServerError = errors.New("internal server")
	ErrMsgValidation          = errors.New("validation error")
)

// Definition Struct (Class)
type UserHandler struct {
	usecase  usecase.UserUsecase
	validate *validator.Validate
}

// Definition Constructor
func NewUserHandler(usecase usecase.UserUsecase, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		usecase:  usecase,
		validate: validate,
	}
}

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
		return
	}

	err = h.validate.Struct(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpresponse.Error(ErrMsgValidation.Error(), err))
		return
	}

	userCreate, err := h.usecase.Create(c, entities.User{
		Name:   req.Name,
		Family: req.Family,
		Email:  req.Email,
		Age:    req.Age,
	})

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrMsgDuplicateUser):
			c.JSON(http.StatusConflict, httpresponse.Error(ErrMsgDuplicateUser.Error(), err))
			return
		case errors.Is(err, usecase.ErrMsgFailedToSaveUser):
			c.JSON(http.StatusInternalServerError, httpresponse.Error(ErrMsgFailedToSaveUser.Error(), err))
			return
		default:
			c.JSON(http.StatusInternalServerError, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
			return
		}
	}

	c.JSON(http.StatusOK, httpresponse.Success("User Created Successfully", response.CreatedUserResponse{
		ID:     userCreate.ID,
		Name:   userCreate.Name,
		Family: userCreate.Family,
		Email:  userCreate.Email,
		Age:    userCreate.Age,
	}))
}
