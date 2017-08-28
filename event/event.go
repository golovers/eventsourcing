package event

import "github.com/lnquy/eventsourcing/model"

type (
	PersonCreated struct {
		model.Person
	}

	PersonUpdated struct {
		model.Person
	}

	PersonDeleted struct {
		Id string
	}
)
