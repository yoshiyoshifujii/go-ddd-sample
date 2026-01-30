package user

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var ErrInvalidUserID = errors.New("invalid user id")

// UserID is a value object.
type UserID string

func NewUserID() (UserID, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return UserID(hex.EncodeToString(b[:])), nil
}

func UserIDFromString(value string) (UserID, error) {
	if value == "" {
		return "", ErrInvalidUserID
	}
	return UserID(value), nil
}

func (id UserID) String() string {
	return string(id)
}

func (id UserID) IsZero() bool {
	return id == ""
}
