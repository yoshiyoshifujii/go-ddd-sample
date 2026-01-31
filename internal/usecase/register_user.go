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

type RegisterUserUsecaseInput struct {
	Name  user.UserName
	Email user.Email
}

type RegisterUserUsecaseOutput struct {
	UserID string
}

func (u *RegisterUser) Execute(ctx context.Context, input RegisterUserUsecaseInput) (*RegisterUserUsecaseOutput, error) {
	if _, err := u.repo.FindByEmail(ctx, input.Email); err == nil {
		return nil, ErrEmailAlreadyUsed
	} else if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	id := user.NewUserID(user.GenerateUserID())

	entity := user.NewUser(id, input.Name, input.Email, time.Now().UTC())

	if err := u.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &RegisterUserUsecaseOutput{UserID: id.String()}, nil
}
