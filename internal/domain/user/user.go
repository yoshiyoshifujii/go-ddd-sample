package user

import "time"

// User is an entity with identity and invariants.
type User struct {
	id        UserID
	name      UserName
	email     Email
	createdAt time.Time
}

func NewUser(id UserID, name UserName, email Email, createdAt time.Time) User {
	if id.IsZero() || name.IsZero() || email.IsZero() {
		panic("invalid user")
	}
	if createdAt.IsZero() {
		panic("invalid user state")
	}
	return User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: createdAt,
	}
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
