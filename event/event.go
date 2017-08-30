package event

const (
	// Event types
	PersonCreatedEvent = "person:created"
	PersonUpdatedEvent = "person:updated"
)

type (
	// Struct for PersonCreated event.
	PersonCreated struct {
		Id   string
		Name string
		Age  int
	}

	// Struct for PersonUpdated event.
	PersonUpdated struct {
		Name string
		Age  int
	}

	// Event represents a state change that has occurred.
	// Event is the root wrapper for all other event types.
	Event struct {
		// The UUID of the event, if not populated the server will create one.
		Id string
		// The event payload. Contains the actual event.
		Data []byte
		// Type of event
		Type string
		// If left blank, this will be populated with the server-time once written to the log.
		Timestamp string
	}
)
