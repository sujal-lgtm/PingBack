package services

import (
	"fmt"
)

type Replayer struct{
	store *Store
	forwarder *Forwarder
}

type Forwarder struct{}

func (f *Forwarder) ForwardEvent(payload string, targetURL string) error {
	return ForwardEvent(payload, targetURL)
}

func NewReplayer(store *Store, forwarder *Forwarder) *Replayer{
	return &Replayer{
		store: store,
		forwarder: forwarder,
	}
}

func (r *Replayer) ReplayEvent (eventID, targetURL string) error{
	event, found := r.store.Get(eventID)
	if !found{
		return fmt.Errorf("Event not found")
	}

	return r.forwarder.ForwardEvent(event.Payload, targetURL)
}