package eventsdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/utils/hash"
)

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
	topic       = "Event"
)

// Event exposes event APIs of MESG.
type Event struct {
	ps         *pubsub.PubSub
	db         database.ServiceDB
	instanceDB database.InstanceDB
}

// New creates a new Event SDK with given options.
func New(ps *pubsub.PubSub, db database.ServiceDB, instanceDB database.InstanceDB) *Event {
	return &Event{
		ps:         ps,
		db:         db,
		instanceDB: instanceDB,
	}
}

// Emit emits a MESG event eventKey with eventData for instanceHash.
func (e *Event) Emit(instanceHash, eventKey string, eventData map[string]interface{}) error {
	i, err := e.instanceDB.Get(instanceHash)
	if err != nil {
		return err
	}
	s, err := e.db.Get(i.ServiceHash)
	if err != nil {
		return err
	}
	ev, err := event.Create(s, i, eventKey, eventData)
	if err != nil {
		return err
	}

	go e.ps.Pub(ev, streamTopic)
	go e.ps.Pub(ev, subTopic(instanceHash))
	return nil
}

// GetStream broadcasts all events.
func (e *Event) GetStream(f *Filter) *Listener {
	l := NewListener(e.ps, streamTopic, f)
	go l.Listen()
	return l
}

// Listen listens events matches with eventFilter on instanceHash.
func (e *Event) Listen(instanceHash string, f *Filter) (*Listener, error) {
	i, err := e.instanceDB.Get(instanceHash)
	if err != nil {
		return nil, err
	}
	s, err := e.db.Get(i.ServiceHash)
	if err != nil {
		return nil, err
	}

	if f.HasKey() {
		if _, err := s.GetEvent(f.Key); err != nil {
			return nil, err
		}
	}

	l := NewListener(e.ps, subTopic(instanceHash), f)
	go l.Listen()
	return l, nil
}

// subTopic returns the topic to listen for events from this service.
func subTopic(serviceHash string) string {
	return hash.Calculate([]string{serviceHash, topic})
}
