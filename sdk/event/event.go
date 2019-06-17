package eventsdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/utils/hash"
)

const (
	topic = "Event"
)

// Event exposes event APIs of MESG.
type Event struct {
	ps *pubsub.PubSub
	db database.ServiceDB
}

// New creates a new Event SDK with given options.
func New(ps *pubsub.PubSub, db database.ServiceDB) *Event {
	return &Event{
		ps: ps,
		db: db,
	}
}

// Emit emits a MESG event eventKey with eventData for service token.
func (e *Event) Emit(token, eventKey string, eventData map[string]interface{}) error {
	s, err := e.db.Get(token)
	if err != nil {
		return err
	}
	ev, err := event.Create(s, eventKey, eventData)
	if err != nil {
		return err
	}

	go e.ps.Pub(ev, subTopic(s.Hash))
	return nil
}

// Listen listens events matches with eventFilter on serviceID.
func (e *Event) Listen(service string, f *Filter) (*Listener, error) {
	s, err := e.db.Get(service)
	if err != nil {
		return nil, err
	}

	if f.HasKey() {
		if _, err := s.GetEvent(f.Key); err != nil {
			return nil, err
		}
	}

	l := NewListener(e.ps, subTopic(s.Hash), f)
	go l.Listen()
	return l, nil
}

// subTopic returns the topic to listen for events from this service.
func subTopic(serviceHash string) string {
	return hash.Calculate([]string{serviceHash, topic})
}
