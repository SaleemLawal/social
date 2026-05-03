package cache

import (
	"context"

	"github.com/saleemlawal/social/internal/store"
)

func NewMockCache() Storage {
	return Storage{
		Users: &MockUsersStore{},
	}
}

type MockUsersStore struct{}

func (m *MockUsersStore) Get(ctx context.Context, id int64) (*store.User, error) {
	return nil, nil
}

func (m *MockUsersStore) Set(ctx context.Context, user *store.User) error {
	return nil
}
