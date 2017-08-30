package event

const (
	PersonCreatedEvent = "person:created"
	PersonUpdatedEvent = "person:updated"
)

type (
	PersonCreated struct {
		Id string
		Name string
		Age  int
	}

	PersonUpdated struct {
		Name string
		Age  int
	}
)

// Event represents a state change that has occurred.
type Event struct {
	// The UUID of the event, if not populated the server will create one.
	Id string
	// The event payload. For JSON requests, the value of this field must be base64-encoded.
	Data []byte
	// Type of event
	Type string
	// If left blank, this will be populated with the server-time once written
	// to the log.
	Timestamp string
}
