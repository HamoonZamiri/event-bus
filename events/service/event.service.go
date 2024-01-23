package service

import (
	"example/event-bus/model"
	"slices"
)

type EventStore struct {
	Subscribers map[string][]string
	Events map[string][]*model.UnknownEvent
}

type EventService interface {
	PostEventType(eventType string) error
	Subscribe(eventType string, event interface{}) error
	PublishEvent(eventType string, event model.UnknownEvent)error
	GetSubscribers(eventType string) ([]string, error)
	GetEventsByType(eventType string) ([]*model.UnknownEvent, error)
	DeleteSubscriber(eventType string, host string) error
	ListEventTypes() ([]string, error)
}

func NewEventStore(eventTypes []string) *EventStore {
	var eventStore = new(EventStore)
	eventStore.Subscribers = make(map[string][]string)
	eventStore.Events = make(map[string][]*model.UnknownEvent)
	for _, eventType := range eventTypes {
		eventStore.Subscribers[eventType] = []string{}
	}
	return eventStore
}

func (es *EventStore) Subscribe(eventType string, event interface{}) error {
	es.Subscribers[eventType] = append(es.Subscribers[eventType], event.(string))
	return nil
}

func (es *EventStore) PublishEvent(eventType string, event model.UnknownEvent) error {
	es.Events[eventType] = append(es.Events[eventType], &event)
	return nil
}

func (es *EventStore) GetSubscribers(eventType string) ([]string, error) {
	return es.Subscribers[eventType], nil
}

func (es *EventStore) GetEventsByType(eventType string) ([]*model.UnknownEvent, error) {
	return es.Events[eventType], nil
}

func (es *EventStore) DeleteSubscriber(eventType string, host string) error {
	var subscribers = es.Subscribers[eventType]
	for i, sub := range subscribers {
		if sub == host {
			subscribers = slices.Delete(subscribers, i, i+1)
			es.Subscribers[eventType] = subscribers
			break
		}
	}
	return nil
}

func (es *EventStore) ListEventTypes() ([]string, error) {
	var eventTypes []string
	for eventType := range es.Events {
		eventTypes = append(eventTypes, eventType)
	}
	return eventTypes, nil
}