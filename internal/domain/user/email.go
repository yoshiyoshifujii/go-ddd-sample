package user

import (
	"net/mail"
	"strings"
)

// Email is a value object.
type Email string

func NewEmail(value string) Email {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		panic("invalid email")
	}
	if _, err := mail.ParseAddress(trimmed); err != nil {
		panic("invalid email")
	}
	return Email(strings.ToLower(trimmed))
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsZero() bool {
	return e == ""
}
