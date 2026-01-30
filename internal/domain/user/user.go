package user

import (
	"errors"
	"time"
)

var (
	ErrInvalidUser      = errors.New("invalid user")
	ErrInvalidUserState = errors.New("invalid user state")
)

// User is an entity with identity and invariants.
type User struct {
	id        UserID
	name      UserName
	email     Email
	createdAt time.Time
}

func NewUser(id UserID, name UserName, email Email, createdAt time.Time) (*User, error) {
	if id.IsZero() || name.IsZero() || email.IsZero() {
		return nil, ErrInvalidUser
	}
	if createdAt.IsZero() {
		return nil, ErrInvalidUserState
	}
	return &User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: createdAt,
	}, nil
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Name() UserName {
	return u.name
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) ChangeName(name UserName) error {
	if name.IsZero() {
		return ErrInvalidUser
	}
	u.name = name
	return nil
}
