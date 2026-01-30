package usecase_test

import (
	"context"
	"errors"
	"testing"

	"yoshiyoshifujii/go-ddd-sample/internal/infrastructure/memory"
	"yoshiyoshifujii/go-ddd-sample/internal/usecase"
)

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()
	uc := usecase.NewRegisterUser(repo)

	out, err := uc.Execute(ctx, usecase.RegisterUserInput{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.UserID == "" {
		t.Fatalf("expected user id")
	}

	_, err = uc.Execute(ctx, usecase.RegisterUserInput{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	if !errors.Is(err, usecase.ErrEmailAlreadyUsed) {
		t.Fatalf("expected ErrEmailAlreadyUsed, got %v", err)
	}
}
