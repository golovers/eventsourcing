package datastore

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"github.com/Sirupsen/logrus"
	"github.com/lnquy/eventsourcing/model"
)

const kind = "tma-eventstore.v1.Person"

type PersonStorage struct {
	client      *datastore.Client
	namespace   string
}

// NewStorage creates a new PersonStorage client.
func NewStorage(ctx context.Context, cfg *Config) (*PersonStorage, error) {
	client, err := datastore.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		logrus.Errorf("datastore: failed to create client: %v", err)
		return nil, err
	}

	return &PersonStorage {
		client:      client,
		namespace:   cfg.Namespace,
	}, nil
}

// Get returns a Person from its resource name.
func (s *PersonStorage) Get(ctx context.Context, name string) (retDc *model.Person, retErr error) {
	retDc = new(model.Person)
	retErr = s.client.Get(ctx, s.key(name), retDc)
	return
}

// Save stores a Person to datastore and returns the saved Person.
// If the Person is not existed in datastore, current Person will be saved (create).
// Otherwise, current Person will override the existing Person in datastore (update).
func (s *PersonStorage) Save(ctx context.Context, dc *model.Person) (retDc *model.Person, retErr error) {
	retDc = new(model.Person)
	_, retErr = s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := s.key(dc.Name)
		curPerson := &model.Person{}

		if err := tx.Get(key, curPerson); err == nil {
			curPerson.ApplyFrom(dc) // PATCH: Override current person fields
		} else if err == datastore.ErrNoSuchEntity {
			curPerson = dc // POST: Create new person
		} else {
			return err
		}
		// Put data connector to datastore
		if _, err := tx.Put(key, curPerson); err != nil {
			return err
		}
		retDc = curPerson
		return nil
	})
	return
}

// Delete deletes a Person from datastore by its resource name.
func (s *PersonStorage) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, s.key(name))
}

// key returns the datastore.Key from the Person resource name.
func (s *PersonStorage) key(id string) *datastore.Key {
	key := datastore.NameKey(kind, id, nil)
	key.Namespace = s.namespace
	return key
}
