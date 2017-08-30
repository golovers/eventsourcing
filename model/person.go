package model

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/lnquy/eventsourcing/event"
	"github.com/lnquy/eventsourcing/store"
	"golang.org/x/net/context"
)

var s = store.CreateMemStore() // Memory storage

// Person represents a person.
type Person struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

// SavePersonEvents stores the event log of the person to memory storage.
func SavePersonEvents(ctx context.Context, event *event.Event, id string) {
	s.Commit(ctx, event, id)
}

// GetPersonAggregate replays all committed event logs to rebuild the Person aggregate.
func GetPersonAggregate(ctx context.Context, p *Person, id string) error {
	events, err := s.Replay(ctx, id)
	if err != nil {
		return err
	}
	for _, c := range events {
		if err := p.applyEvent(c); err != nil {
			return err
		}
	}
	return nil
}

// applyEvent applies an event to the person aggregate.
func (p *Person) applyEvent(e *event.Event) error {
	switch e.Type {
	case event.PersonCreatedEvent:
		pc := &event.PersonCreated{}
		err := json.Unmarshal(e.Data, pc)
		if err != nil {
			logrus.Errorf("model: failed to decode PersonCreated event: %v", err)
			return err
		}
		p.apply(pc.Id, pc.Name, pc.Age)
	case event.PersonUpdatedEvent:
		pu := &event.PersonUpdated{}
		err := json.Unmarshal(e.Data, pu)
		if err != nil {
			logrus.Errorf("model: failed to decode PersonUpdated event: %v", err)
			return err
		}
		p.apply("", pu.Name, pu.Age)
	}
	return nil
}

func (p *Person) apply(id, name string, age int) {
	if id != "" {
		p.Id = id
	}
	if name != "" {
		p.Name = name
	}
	if age >= 0 {
		p.Age = age
	}
}
