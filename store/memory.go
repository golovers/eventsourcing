// Package store provides an in-memory database to stores event.Events.
package store

import (
	"context"
	"sync"
	"fmt"
	"github.com/lnquy/eventsourcing/event"
)

// Store is the data structure for the in-memory database.
type Store struct {
	eventsMux *sync.RWMutex
	events    map[string][]*event.Event
}

// CreateMemStore returns a in-memory database Store.
func CreateMemStore() *Store {
	return &Store{
		eventsMux: &sync.RWMutex{},
		events:    make(map[string][]*event.Event),
	}
}

// Commit saves an event to its corresponding id in the in-memory database.
func (s *Store) Commit(ctx context.Context, e *event.Event, id string) {
	s.eventsMux.Lock()
	s.events[id] = append(s.events[id], e)
	s.eventsMux.Unlock()
}

// Replay returns the list of all committed events of the id.
func (s *Store) Replay(ctx context.Context, id string) ([]*event.Event, error) {
	s.eventsMux.RLock()
	defer s.eventsMux.RUnlock()
	events, exists := s.events[id]
	if !exists {
		return nil, fmt.Errorf("mem: failed to get event logs for %s id", id)
	}
	return events, nil
}
