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
	ErrMsgDuplicateUser    = errors.New("duplicate user")
	ErrMsgFailedToSaveUser = errors.New("failed to save user")
)

// Definition Interface (Rules)
type UserRepository interface {
	Save(ctx context.Context, user entities.User) (entities.User, error)
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
