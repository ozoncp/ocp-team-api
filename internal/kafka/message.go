package kafka

type Event int

const (
	Create Event = iota + 1
	Update
	Delete
)

func (e Event) String() string {
	return [...]string{"Create", "Update", "Delete"}[e-1]
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
