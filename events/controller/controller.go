package controller

import "example/event-bus/service"

type Controller struct {
	eventStore *service.EventStore
}

func NewController(eventStore *service.EventStore) *Controller {
	return &Controller{eventStore: eventStore}
}