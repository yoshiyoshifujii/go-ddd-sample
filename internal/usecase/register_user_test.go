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
	tests := []struct {
		name        string
		firstInput  usecase.RegisterUserUsecaseInput
		secondInput usecase.RegisterUserUsecaseInput
		wantSameID  bool
	}{
		{
			name: "same email returns same user id",
			firstInput: usecase.RegisterUserUsecaseInput{
				Name:  user.NewUserName("Alice"),
				Email: user.NewEmail("alice@example.com"),
			},
			secondInput: usecase.RegisterUserUsecaseInput{
				Name:  user.NewUserName("Alice"),
				Email: user.NewEmail("alice@example.com"),
			},
			wantSameID: true,
		},
		{
			name: "different email returns different user id",
			firstInput: usecase.RegisterUserUsecaseInput{
				Name:  user.NewUserName("Alice"),
				Email: user.NewEmail("alice@example.com"),
			},
			secondInput: usecase.RegisterUserUsecaseInput{
				Name:  user.NewUserName("Bob"),
				Email: user.NewEmail("bob@example.com"),
			},
			wantSameID: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()
			repo := memory.NewUserRepository()
			idGenerator := memory.NewUserIDGenerator()

			sut := usecase.NewRegisterUser(repo, idGenerator)

			// when
			out, err := sut.Execute(ctx, tt.firstInput)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, out)
			assert.NotEmpty(t, out.UserID)

			// when
			out2, err := sut.Execute(ctx, tt.secondInput)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, out2)
			if tt.wantSameID {
				assert.Equal(t, out.UserID, out2.UserID)
			} else {
				assert.NotEqual(t, out.UserID, out2.UserID)
			}
		})
	}
}
