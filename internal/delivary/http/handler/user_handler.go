package handler

import (
	"errors"
	"net/http"
	"strconv"

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
	ErrMsgFailedToUpdateUser  = errors.New("failed to update user")
	ErrMsgUserNotFound        = errors.New("user not found")
	ErrMsgInvalidId           = errors.New("invalid id")
	ErrMsgValidation          = errors.New("validation error")

	SuccessMsgCreatedUser   = "User Created Successfully"
	SuccessMsgFoundUserById = "found user successfully"
	SuccessMsgUpdatedUser   = "User Updated Successfully"
	SuccessMsgFoundAllUser  = "found users successfully"
	SuccessMsgDeletedUser   = "User Deleted Successfully"
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

	c.JSON(http.StatusOK, httpresponse.Success(SuccessMsgCreatedUser, response.CreatedUserResponse{
		ID:     userCreate.ID,
		Name:   userCreate.Name,
		Family: userCreate.Family,
		Email:  userCreate.Email,
		Age:    userCreate.Age,
	}))
}

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (h *UserHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpresponse.Error(ErrMsgInvalidId.Error(), err))
		return
	}

	user, err := h.usecase.FindByID(c, id)
	if err != nil {
		if errors.Is(err, usecase.ErrMsgUserNotFound) {
			c.JSON(http.StatusNotFound, httpresponse.Error(ErrMsgUserNotFound.Error(), err))
			return
		}

		c.JSON(http.StatusInternalServerError, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
		return
	}

	c.JSON(http.StatusOK, httpresponse.Success(SuccessMsgFoundUserById, response.FindUserByIDResponse{
		ID:     user.ID,
		Name:   user.Name,
		Family: user.Family,
		Email:  user.Email,
		Age:    user.Age,
	}))
}

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpresponse.Error(ErrMsgInvalidId.Error(), err))
		return
	}

	var req request.UpdateUserRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
		return
	}

	err = h.validate.Struct(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpresponse.Error(ErrMsgValidation.Error(), err))
		return
	}

	// check exists user for update
	existingUser, err := h.usecase.FindByID(c, id)
	if err != nil {
		if errors.Is(err, usecase.ErrMsgUserNotFound) {
			c.JSON(http.StatusNotFound, httpresponse.Error(ErrMsgUserNotFound.Error(), err))
			return
		}

		c.JSON(http.StatusInternalServerError, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
		return
	}

	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	if req.Family != nil {
		existingUser.Family = *req.Family
	}
	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.Age != nil {
		existingUser.Age = *req.Age
	}

	updatedUser, err := h.usecase.Update(c, existingUser)
	if err != nil {
		if errors.Is(err, usecase.ErrMsgDuplicateUser) {
			c.JSON(http.StatusConflict, httpresponse.Error(ErrMsgDuplicateUser.Error(), err))
			return
		}

		if errors.Is(err, usecase.ErrMsgUserNotFound) {
			c.JSON(http.StatusNotFound, httpresponse.Error(ErrMsgUserNotFound.Error(), err))
			return
		}

		c.JSON(http.StatusInternalServerError, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
		return
	}

	c.JSON(http.StatusOK, httpresponse.Success(SuccessMsgUpdatedUser, response.UpdatedUserResponse{
		ID:     updatedUser.ID,
		Name:   updatedUser.Name,
		Family: updatedUser.Family,
		Email:  updatedUser.Email,
		Age:    updatedUser.Age,
	}))

}

func (h *UserHandler) FindAll(c *gin.Context) {
	users, err := h.usecase.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
		return
	}

	var res []response.FindAllUserResponse
	for _, user := range users {
		res = append(res, response.FindAllUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Family: user.Family,
			Email:  user.Email,
			Age:    user.Age,
		})
	}

	c.JSON(http.StatusOK, httpresponse.Success(SuccessMsgFoundAllUser, res))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpresponse.Error(ErrMsgInvalidId.Error(), err))
		return
	}

	err = h.usecase.Delete(c, id)
	if err != nil {
		if errors.Is(err, usecase.ErrMsgUserNotFound) {
			c.JSON(http.StatusNotFound, httpresponse.Error(ErrMsgUserNotFound.Error(), err))
			return
		}

		c.JSON(http.StatusInternalServerError, httpresponse.Error(ErrMsgInternalServerError.Error(), err))
		return
	}

	c.JSON(http.StatusOK, httpresponse.Success(SuccessMsgDeletedUser, response.DeletedUserResponse{
		ID: id,
	}))
}
