package memory

import (
	"context"
	"sync"

	"yoshiyoshifujii/go-ddd-sample/internal/domain/user"
)

// UserIDGenerator is an in-memory implementation of generator.UserIDGenerator.
type UserIDGenerator struct {
	mu      sync.RWMutex
	byEmail map[user.Email]user.UserID
}

func NewUserIDGenerator() *UserIDGenerator {
	return &UserIDGenerator{
		byEmail: make(map[user.Email]user.UserID),
	}
}

func (g *UserIDGenerator) Generate(ctx context.Context, email user.Email) (user.UserID, error) {
	_ = ctx

	g.mu.RLock()
	id, ok := g.byEmail[email]
	g.mu.RUnlock()
	if ok {
		return id, nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if id, ok := g.byEmail[email]; ok {
		return id, nil
	}

	id = user.NewUserID(user.GenerateUserID())
	g.byEmail[email] = id
	return id, nil
}
