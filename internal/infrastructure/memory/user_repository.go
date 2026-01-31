package memory

import (
	"context"
	"sync"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
	"yoshiyoshifujii/go-ddd-sample/internal/repository"
)

// UserRepository is an in-memory implementation for tests or local usage.
type UserRepository struct {
	mu   sync.RWMutex
	byID map[user.UserID]user.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID: make(map[user.UserID]user.User),
	}
}

func (r *UserRepository) Save(ctx context.Context, entity user.User) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[entity.ID()] = entity
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (user.User, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, ok := r.byID[id]
	if !ok {
		return user.User{}, repository.ErrUserNotFound
	}
	return entity, nil
}
