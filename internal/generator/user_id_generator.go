package generator

import (
	"context"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
)

// UserIDGenerator generates stable user IDs by email.
type UserIDGenerator interface {
	Generate(ctx context.Context, email user.Email) (user.UserID, error)
}
