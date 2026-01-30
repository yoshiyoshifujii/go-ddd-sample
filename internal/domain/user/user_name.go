package user

import (
	"errors"
	"strings"
)

var ErrInvalidUserName = errors.New("invalid user name")

// UserName is a value object.
type UserName string

func UserNameFromString(value string) (UserName, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || len(trimmed) > 50 {
		return "", ErrInvalidUserName
	}
	return UserName(trimmed), nil
}

func (name UserName) String() string {
	return string(name)
}

func (name UserName) IsZero() bool {
	return name == ""
}
