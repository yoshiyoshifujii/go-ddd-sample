package user

import "strings"

// UserName is a value object.
type UserName string

func NewUserName(value string) UserName {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || len(trimmed) > 50 {
		panic("invalid user name")
	}
	return UserName(trimmed)
}

func (name UserName) String() string {
	return string(name)
}

func (name UserName) IsZero() bool {
	return name == ""
}
