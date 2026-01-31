package usecase

import (
	"context"
	"errors"
	"time"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
	"yoshiyoshifujii/go-ddd-sample/internal/generator"
	"yoshiyoshifujii/go-ddd-sample/internal/repository"
)

// RegisterUser registers a new user.
type RegisterUser struct {
	repo        repository.UserRepository
	idGenerator generator.UserIDGenerator
}

func NewRegisterUser(repo repository.UserRepository, idGenerator generator.UserIDGenerator) *RegisterUser {
	return &RegisterUser{repo: repo, idGenerator: idGenerator}
}

type RegisterUserUsecaseInput struct {
	Name  user.UserName
	Email user.Email
}

type RegisterUserUsecaseOutput struct {
	UserID string
}

func (u *RegisterUser) Execute(ctx context.Context, input RegisterUserUsecaseInput) (*RegisterUserUsecaseOutput, error) {
	id, err := u.idGenerator.Generate(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	entity, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
	} else {
		return &RegisterUserUsecaseOutput{UserID: entity.ID().String()}, nil
	}

	entity = user.NewUser(id, input.Name, input.Email, time.Now().UTC())

	if err := u.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &RegisterUserUsecaseOutput{UserID: id.String()}, nil
}
