package user

import (
	"crypto/rand"
	"encoding/hex"
)

// UserID is a value object.
type UserID string

func NewUserID(value string) UserID {
	if value == "" {
		panic("invalid user id")
	}
	return UserID(value)
}

func GenerateUserID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}

func (id UserID) String() string {
	return string(id)
}

func (id UserID) IsZero() bool {
	return id == ""
}
