package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alirezaghasemi/user-manager/internal/entities"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Definition Error Message
var (
	ErrMsgDuplicateUser       = errors.New("duplicate user")
	ErrMsgFailedToSaveUser    = errors.New("failed to save user")
	ErrMsgFailedToUpdateUser  = errors.New("failed to update user")
	ErrMsgUserNotFound        = errors.New("user not found")
	ErrMsgInternalServerError = errors.New("internal server")
)

// Definition Interface (Rules)
type UserRepository interface {
	Save(ctx context.Context, user entities.User) (entities.User, error)
	FindByID(ctx context.Context, id uint64) (entities.User, error)
	Update(ctx context.Context, user entities.User) (entities.User, error)
	FindAll(ctx context.Context) ([]entities.User, error)
	Delete(ctx context.Context, id uint64) error
}

// Definition Struct (Class)
type userRepository struct {
	db *gorm.DB
}

// Definition Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Definition Implement Methods (Save, Update, Delete, FindByID, FindAll)
func (r *userRepository) Save(ctx context.Context, user entities.User) (entities.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		// below check error for duplicate record error
		// var pgErr *pq.Error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgDuplicateUser, err)
		}
		return entities.User{}, fmt.Errorf("%w:%w", ErrMsgFailedToSaveUser, err)
	}

	return user, nil
}

// Definition Implement Methods (Save, Update, Delete, FindByID, FindAll)
func (r *userRepository) FindByID(ctx context.Context, id uint64) (entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgUserNotFound, err)
		}

		return entities.User{}, fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	return user, nil
}

// Definition Implement Methods (Save, Update, Delete, FindByID, FindAll)
func (r *userRepository) Update(ctx context.Context, user entities.User) (entities.User, error) {
	tx := r.db.WithContext(ctx).Model(&user).Where("id = ?", user.ID).Updates(user)

	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) && pgErr.Code == "23505" {
			return entities.User{}, fmt.Errorf("%w:%w", ErrMsgDuplicateUser, tx.Error)
		}

		return entities.User{}, fmt.Errorf("%w:%w", ErrMsgFailedToUpdateUser, tx.Error)
	}

	if tx.RowsAffected == 0 {
		return entities.User{}, fmt.Errorf("%w:%w", ErrMsgUserNotFound, tx.Error)
	}

	return user, nil
}

// Definition Implement Methods (Save, Update, Delete, FindByID, FindAll)
func (r *userRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return []entities.User{}, fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	return users, nil
}

// Definition Implement Methods (Save, Update, Delete, FindByID, FindAll)
func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	// first check record is exists
	var user entities.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w:%w", ErrMsgUserNotFound, err)
		}

		return fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	// delete section
	err = r.db.WithContext(ctx).Delete(&user).Error
	if err != nil {
		return fmt.Errorf("%w:%w", ErrMsgInternalServerError, err)
	}

	return nil
}
