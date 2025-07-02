package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alirezaghasemi/user-manager/internal/entities"
	"github.com/alirezaghasemi/user-manager/internal/repository"
	"github.com/go-playground/validator/v10"
)

// Definition Error Message
var (
	ErrMsgDuplicateUser       = errors.New("duplicate user")
	ErrMsgFailedToSaveUser    = errors.New("failed to save user")
	ErrMsgInternalServerError = errors.New("internal server")
)

// Definition Interface (Rules)
type UserUsecase interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
}

// Definition Struct (Class)
type userUsecase struct {
	repo     repository.UserRepository
	validate *validator.Validate
}

// Definition Constructor
func NewUserUsecase(repo repository.UserRepository, validate *validator.Validate) UserUsecase {
	return &userUsecase{
		repo:     repo,
		validate: validate,
	}
}

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (u *userUsecase) Create(ctx context.Context, user entities.User) (entities.User, error) {
	userSaved, err := u.repo.Save(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrMsgDuplicateUser) {
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgDuplicateUser, err)
		}

		if errors.Is(err, repository.ErrMsgFailedToSaveUser) {
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgFailedToSaveUser, err)
		}

		return entities.User{}, fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	return userSaved, nil
}
