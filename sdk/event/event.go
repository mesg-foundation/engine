package eventsdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/utils/hash"
	"github.com/mr-tron/base58"
)

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
	topic       = "Event"
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
func (e *Event) Emit(serviceHash []byte, eventKey string, eventData map[string]interface{}) error {
	s, err := e.db.Get(serviceHash)
	if err != nil {
		return err
	}
	ev, err := event.Create(s, eventKey, eventData)
	if err != nil {
		return err
	}

	go e.ps.Pub(ev, streamTopic)
	go e.ps.Pub(ev, subTopic(s.Hash))
	return nil
}

// GetStream broadcasts all events.
func (e *Event) GetStream(f *Filter) *Listener {
	l := NewListener(e.ps, streamTopic, f)
	go l.Listen()
	return l
}

// Listen listens events matches with eventFilter on serviceID.
func (e *Event) Listen(serviceHash []byte, f *Filter) (*Listener, error) {
	s, err := e.db.Get(serviceHash)
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
func subTopic(serviceHash []byte) string {
	return hash.Calculate([]string{base58.Encode(serviceHash), topic})
}
