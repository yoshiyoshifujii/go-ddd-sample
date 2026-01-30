package user

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository defines persistence for User aggregate.
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
}
