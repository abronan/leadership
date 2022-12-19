package leadership

import (
	"abronan/leadership/mockstore"
	"context"
	"testing"

	"github.com/kvtools/valkeyrie/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFollower(t *testing.T) {
	kv, err := mockstore.New([]string{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	mockStore := kv.(*mockstore.Mock)

	kvCh := make(chan *store.KVPair)
	var mockKVCh <-chan *store.KVPair = kvCh

	ctx, _ := context.WithCancel(context.Background())

	mockStore.On("Watch", ctx, mock.Anything, mock.AnythingOfType("*store.ReadOptions")).Return(mockKVCh, nil)

	follower := NewFollower(kv, "test_key")

	// Calling stop when not following an election should not panic and should
	// result in no side effects observed when starting to follow an election.
	follower.Stop()

	leaderCh, errCh := follower.FollowElection()

	// Simulate leader updates
	go func() {
		kvCh <- &store.KVPair{Key: "test_key", Value: []byte("leader1")}
		kvCh <- &store.KVPair{Key: "test_key", Value: []byte("leader1")}
		kvCh <- &store.KVPair{Key: "test_key", Value: []byte("leader2")}
		kvCh <- &store.KVPair{Key: "test_key", Value: []byte("leader1")}
	}()

	// We shouldn't see duplicate events.
	assert.Equal(t, <-leaderCh, "leader1")
	assert.Equal(t, <-leaderCh, "leader2")
	assert.Equal(t, <-leaderCh, "leader1")
	assert.Equal(t, follower.Leader(), "leader1")

	// Once stopped, iteration over the leader channel should stop.
	follower.Stop()
	close(kvCh)

	// Assert that we receive an error from the error chan to deal with the failover
	err, open := <-errCh
	assert.True(t, open)
	assert.Error(t, err)
	assert.Equal(t, err, ErrStoreUnavailable)

	// Ensure that the chan is closed
	_, open = <-leaderCh
	assert.False(t, open)

	mockStore.AssertExpectations(t)
}
