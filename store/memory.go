// Package store provides an in-memory implementation of a storage.Store.
package store

import (
	"context"
	"sync"
	"fmt"
	"github.com/lnquy/eventsourcing/event"
)

type Store struct {
	eventsMux *sync.RWMutex
	events   map[string][]*event.Event
}

func CreateMemStore() *Store {
	return &Store{
		eventsMux: &sync.RWMutex{},
		events:   make(map[string][]*event.Event),
	}
}

// Commit saves an provided commit at the specified commit sequence.
func (s *Store) Commit(ctx context.Context, e *event.Event, id string) {
	s.eventsMux.Lock()
	s.events[id] = append(s.events[id], e)
	s.eventsMux.Unlock()
}

// Replay all commits starting from the provided commit sequence.
func (s *Store) Replay(ctx context.Context, id string) ([]*event.Event, error) {
	s.eventsMux.RLock()
	defer s.eventsMux.RUnlock()
	events, exists := s.events[id]
	if !exists {
		return nil, fmt.Errorf("mem: failed to get event logs for %s id", id)
	}
	return events, nil
}
