package user

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrInvalidEmail = errors.New("invalid email")

// Email is a value object.
type Email string

func EmailFromString(value string) (Email, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(trimmed); err != nil {
		return "", ErrInvalidEmail
	}
	return Email(strings.ToLower(trimmed)), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsZero() bool {
	return e == ""
}
