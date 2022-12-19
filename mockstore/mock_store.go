// Package mockstore Mocks all valkeyrie Store functions using testify.Mock.
package mockstore

import (
	"context"

	"github.com/kvtools/valkeyrie"
	"github.com/kvtools/valkeyrie/store"
	"github.com/stretchr/testify/mock"
)

// Mock store. Mocks all valkeyrie Store functions using testify.Mock.
type Mock struct {
	mock.Mock

	// Endpoints passed to InitializeMock.
	Endpoints []string

	// Options passed to InitializeMock.
	Options *valkeyrie.Config
}

// New creates a Mock store.
func New(endpoints []string, options *valkeyrie.Config) (store.Store, error) {
	s := &Mock{}
	s.Endpoints = endpoints
	s.Options = options
	return s, nil
}

// Put mock.
func (s *Mock) Put(ctx context.Context, key string, value []byte, opts *store.WriteOptions) error {
	args := s.Mock.Called(ctx, key, value, opts)
	return args.Error(0)
}

// Get mock.
func (s *Mock) Get(ctx context.Context, key string, opts *store.ReadOptions) (*store.KVPair, error) {
	args := s.Mock.Called(ctx, key, opts)
	return args.Get(0).(*store.KVPair), args.Error(1)
}

// Delete mock.
func (s *Mock) Delete(ctx context.Context, key string) error {
	args := s.Mock.Called(ctx, key)
	return args.Error(0)
}

// Exists mock.
func (s *Mock) Exists(ctx context.Context, key string, opts *store.ReadOptions) (bool, error) {
	args := s.Mock.Called(ctx, key, opts)
	return args.Bool(0), args.Error(1)
}

// Watch mock.
func (s *Mock) Watch(ctx context.Context, key string, opts *store.ReadOptions) (<-chan *store.KVPair, error) {
	args := s.Mock.Called(ctx, key, opts)
	return args.Get(0).(<-chan *store.KVPair), args.Error(1)
}

// WatchTree mock.
func (s *Mock) WatchTree(ctx context.Context, prefix string, opts *store.ReadOptions) (<-chan []*store.KVPair, error) {
	args := s.Mock.Called(ctx, prefix, opts)
	return args.Get(0).(chan []*store.KVPair), args.Error(1)
}

// NewLock mock.
func (s *Mock) NewLock(ctx context.Context, key string, options *store.LockOptions) (store.Locker, error) {
	args := s.Mock.Called(ctx, key, options)
	return args.Get(0).(store.Locker), args.Error(1)
}

// List mock.
func (s *Mock) List(ctx context.Context, prefix string, opts *store.ReadOptions) ([]*store.KVPair, error) {
	args := s.Mock.Called(ctx, prefix, opts)
	return args.Get(0).([]*store.KVPair), args.Error(1)
}

// DeleteTree mock.
func (s *Mock) DeleteTree(ctx context.Context, prefix string) error {
	args := s.Mock.Called(ctx, prefix)
	return args.Error(0)
}

// AtomicPut mock.
func (s *Mock) AtomicPut(ctx context.Context, key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) {
	args := s.Mock.Called(ctx, key, value, previous, opts)
	return args.Bool(0), args.Get(1).(*store.KVPair), args.Error(2)
}

// AtomicDelete mock.
func (s *Mock) AtomicDelete(ctx context.Context, key string, previous *store.KVPair) (bool, error) {
	args := s.Mock.Called(ctx, key, previous)
	return args.Bool(0), args.Error(1)
}

// Lock mock implementation of Locker.
type Lock struct {
	mock.Mock
}

// Lock mock.
func (l *Lock) Lock(ctx context.Context) (<-chan struct{}, error) {
	args := l.Mock.Called(ctx)
	return args.Get(0).(<-chan struct{}), args.Error(1)
}

// Unlock mock.
func (l *Lock) Unlock(ctx context.Context) error {
	args := l.Mock.Called(ctx)
	return args.Error(0)
}

// Close mock.
func (s *Mock) Close() error {
	args := s.Mock.Called()
	return args.Error(0)
}
