package leadership

import (
	"context"
	"errors"
	"sync"

	"github.com/kvtools/valkeyrie/store"
)

var ErrStoreUnavailable = errors.New("leader election: watch leader channel closed, the store may be unavailable or the watch has been stopped by the caller")

// Follower can follow an election in real-time and push notifications whenever
// there is a change in leadership.
type Follower struct {
	client store.Store
	key    string

	lock       sync.Mutex
	leader     string
	leaderCh   chan string
	cancelFunc func()
	errCh      chan error
}

// NewFollower creates a new follower.
func NewFollower(client store.Store, key string) *Follower {
	return &Follower{
		client: client,
		key:    key,
	}
}

// Leader returns the current leader.
func (f *Follower) Leader() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.leader
}

// FollowElection starts monitoring the election.
func (f *Follower) FollowElection() (<-chan string, <-chan error) {
	f.leaderCh = make(chan string)
	f.errCh = make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	f.cancelFunc = cancel

	go f.follow(ctx)

	return f.leaderCh, f.errCh
}

// Stop stops monitoring an election.
// Calling stop when not following an election results in no effect.
func (f *Follower) Stop() {
	if f.cancelFunc != nil {
		f.cancelFunc()
	}
}

func (f *Follower) follow(ctx context.Context) {
	defer close(f.leaderCh)
	defer close(f.errCh)
	defer f.cancelFunc()

	ch, err := f.client.Watch(ctx, f.key, nil)
	if err != nil {
		f.errCh <- err
		return
	}

	f.leader = ""
	for kv := range ch {
		if kv == nil {
			break
		}

		curr := string(kv.Value)

		f.lock.Lock()
		if curr == f.leader {
			f.lock.Unlock()
			continue
		}
		f.leader = curr
		f.lock.Unlock()

		f.leaderCh <- f.leader
	}

	// Channel closed, we return an error
	f.errCh <- ErrStoreUnavailable
}
