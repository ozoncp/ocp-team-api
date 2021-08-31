package kafka

// Event is the type of action happened: Create, Update, Delete.
type Event int

const (
	Create Event = iota + 1
	Update
	Delete
)

var eventMapper = map[Event]string{
	Create: "Create",
	Update: "Update",
	Delete: "Delete",
}

// String is the method for converting Event type to corresponding string.
func (e Event) String() string {
	if value, ok := eventMapper[e]; ok {
		return value
	}

	return "Unknown"
}

// NewMessage is the constructor method for Message struct.
func NewMessage(Id uint64, Event Event) Message {
	return Message{
		Id:    Id,
		Event: Event.String(),
	}
}

// Message is the struct that representing message to be sent to broker.
type Message struct {
	Id    uint64 `json:"id"`
	Event string `json:"event"`
}
