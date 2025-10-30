package player

import (
	"sync"
)

type queue struct {
	mu        sync.RWMutex
	tracklist []*track
}

func newQueue() *queue {
	return &queue{}
}

func (q *queue) isEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.tracklist) > 0
}

func (q *queue) dequeue() (*track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var t *track

	if len(q.tracklist) == 0 {
		return nil, false
	}

	t = q.tracklist[0]
	q.tracklist = q.tracklist[1:]

	return t, true
}

func (q *queue) enqueue(t *track) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tracklist = append(q.tracklist, t)
}

func (q *queue) clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tracklist = []*track{}
}
