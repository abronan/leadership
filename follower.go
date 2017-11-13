package leadership

import (
	"errors"
	"sync"

	"github.com/abronan/valkeyrie/store"
)

// Follower can follow an election in real-time and push notifications whenever
// there is a change in leadership.
type Follower struct {
	client store.Store
	key    string

	lock     sync.Mutex
	leader   string
	leaderCh chan string
	stopCh   chan struct{}
	errCh    chan error
}

// NewFollower creates a new follower.
func NewFollower(client store.Store, key string) *Follower {
	return &Follower{
		client: client,
		key:    key,
		stopCh: make(chan struct{}),
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

	go f.follow()

	return f.leaderCh, f.errCh
}

// Stop stops monitoring an election.
func (f *Follower) Stop() {
	close(f.stopCh)
}

func (f *Follower) follow() {
	defer close(f.leaderCh)
	defer close(f.errCh)

	ch, err := f.client.Watch(f.key, f.stopCh, nil)
	if err != nil {
		f.errCh <- err
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
	f.errCh <- errors.New("leader Election: watch leader channel closed, the store may be unavailable")
}
