package kafka

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

func (e Event) String() string {
	if value, ok := eventMapper[e]; ok {
		return value
	}

	return "Unknown"
}

func NewMessage(Id uint64, Event Event) Message {
	return Message{
		Id:    Id,
		Event: Event.String(),
	}
}

type Message struct {
	Id    uint64 `json:"id"`
	Event string `json:"event"`
}
