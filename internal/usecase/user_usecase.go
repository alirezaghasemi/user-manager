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
	ErrMsgFailedToUpdateUser  = errors.New("failed to update user")
	ErrMsgUserNotFound        = errors.New("user not found")
)

// Definition Interface (Rules)
type UserUsecase interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	FindByID(ctx context.Context, id uint64) (entities.User, error)
	Update(ctx context.Context, user entities.User) (entities.User, error)
	FindAll(ctx context.Context) ([]entities.User, error)
	Delete(ctx context.Context, id uint64) error
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

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (u *userUsecase) FindByID(ctx context.Context, id uint64) (entities.User, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrMsgUserNotFound) {
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgUserNotFound, err)
		}

		return entities.User{}, fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	return user, nil
}

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (u *userUsecase) Update(ctx context.Context, user entities.User) (entities.User, error) {
	userUpdated, err := u.repo.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMsgDuplicateUser):
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgDuplicateUser, err)
		case errors.Is(err, repository.ErrMsgFailedToUpdateUser):
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgFailedToUpdateUser, err)
		case errors.Is(err, repository.ErrMsgUserNotFound):
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgUserNotFound, err)
		default:
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
		}
	}

	return userUpdated, nil
}

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (u *userUsecase) FindAll(ctx context.Context) ([]entities.User, error) {
	users, err := u.repo.FindAll(ctx)
	if err != nil {
		return []entities.User{}, fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	return users, nil
}

// Definition Implement Methods (Create, Update, Delete, FindByID, FindAll)
func (u *userUsecase) Delete(ctx context.Context, id uint64) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrMsgUserNotFound) {
			return fmt.Errorf("%w:%w", ErrMsgUserNotFound, err)
		}

		return fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	return nil
}
