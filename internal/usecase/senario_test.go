package usecase_test

import (
	"context"
	"testing"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
	"yoshiyoshifujii/go-ddd-sample/internal/infrastructure/memory"
	"yoshiyoshifujii/go-ddd-sample/internal/usecase"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	// given
	ctx := context.Background()
	repo := memory.NewUserRepository()

	sut := usecase.NewRegisterUser(repo)

	input := usecase.RegisterUserUsecaseInput{
		Name:  user.NewUserName("Alice"),
		Email: user.NewEmail("alice@example.com"),
	}

	// when
	out, err := sut.Execute(ctx, input)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotEmpty(t, out.UserID)

	// when
	_, err = sut.Execute(ctx, input)

	// then
	assert.ErrorIs(t, err, usecase.ErrEmailAlreadyUsed)
}
