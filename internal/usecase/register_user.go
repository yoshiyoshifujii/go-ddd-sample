package usecase

import (
	"context"
	"errors"
	"time"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
	"yoshiyoshifujii/go-ddd-sample/internal/repository"
)

var ErrEmailAlreadyUsed = errors.New("email already used")

// RegisterUser registers a new user.
type RegisterUser struct {
	repo repository.UserRepository
}

func NewRegisterUser(repo repository.UserRepository) *RegisterUser {
	return &RegisterUser{repo: repo}
}

type RegisterUserInput struct {
	Name  string
	Email string
}

type RegisterUserOutput struct {
	UserID string
}

func (u *RegisterUser) Execute(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
	name, err := user.UserNameFromString(input.Name)
	if err != nil {
		return RegisterUserOutput{}, err
	}
	email, err := user.EmailFromString(input.Email)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	if _, err := u.repo.FindByEmail(ctx, email); err == nil {
		return RegisterUserOutput{}, ErrEmailAlreadyUsed
	} else if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return RegisterUserOutput{}, err
	}

	id, err := user.NewUserID()
	if err != nil {
		return RegisterUserOutput{}, err
	}

	entity, err := user.NewUser(id, name, email, time.Now().UTC())
	if err != nil {
		return RegisterUserOutput{}, err
	}

	if err := u.repo.Save(ctx, *entity); err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{UserID: id.String()}, nil
}
