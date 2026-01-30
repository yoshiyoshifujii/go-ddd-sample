package memory

import (
	"context"
	"sync"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
)

// UserRepository is an in-memory implementation for tests or local usage.
type UserRepository struct {
	mu      sync.RWMutex
	byID    map[user.UserID]*user.User
	byEmail map[user.Email]*user.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:    make(map[user.UserID]*user.User),
		byEmail: make(map[user.Email]*user.User),
	}
}

func (r *UserRepository) Save(ctx context.Context, entity *user.User) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[entity.ID()] = entity
	r.byEmail[entity.Email()] = entity
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, ok := r.byID[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return entity, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, ok := r.byEmail[email]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return entity, nil
}
