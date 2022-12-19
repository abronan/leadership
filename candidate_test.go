package leadership

import (
	"context"
	"testing"
	"time"

	"github.com/abronan/leadership/mockstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCandidate(t *testing.T) {
	kv, err := mockstore.New([]string{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	mockStore := kv.(*mockstore.Mock)
	mockLock := &mockstore.Lock{}
	mockStore.On("NewLock", context.Background(), mock.Anything, mock.AnythingOfType("*store.LockOptions")).Return(mockLock, nil)

	// Lock and unlock always succeeds.
	lostCh := make(chan struct{})
	var mockLostCh <-chan struct{} = lostCh
	mockLock.On("Lock", context.Background()).Return(mockLostCh, nil)
	mockLock.On("Unlock", context.Background()).Return(nil)

	candidate := NewCandidate(kv, "test_key", "test_node", 0)

	// Calling stop when not running for the election should not panic and should
	// result in no side effects observed when starting to run for the election.
	candidate.Stop()

	electedCh, _ := candidate.RunForElection()

	// Should issue a false upon start, no matter what.
	assert.False(t, <-electedCh)

	// Since the lock always succeeeds, we should get elected.
	assert.True(t, <-electedCh)
	assert.True(t, candidate.IsLeader())

	// Signaling a lost lock should get us de-elected...
	close(lostCh)
	assert.False(t, <-electedCh)

	// And we should attempt to get re-elected again.
	assert.True(t, <-electedCh)

	// When we resign, unlock will get called, we'll be notified of the
	// de-election and we'll try to get the lock again.
	go candidate.Resign()
	assert.False(t, <-electedCh)
	assert.True(t, <-electedCh)

	candidate.Stop()

	// Ensure that the chan closes after some time
	for {
		select {
		case _, open := <-electedCh:
			if !open {
				mockStore.AssertExpectations(t)
				return
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("electedCh not closed correctly")
		}
	}
}
