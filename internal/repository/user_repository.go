package repository

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository defines persistence for User aggregate.
type UserRepository interface {
	Save(ctx context.Context, user user.User) error
	FindByID(ctx context.Context, id user.UserID) (user.User, error)
	FindByEmail(ctx context.Context, email user.Email) (user.User, error)
}
