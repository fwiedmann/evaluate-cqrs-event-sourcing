package infra

import (
	"encoding/json"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user-command/domain"
	"github.com/nats-io/nats.go"
)

type UserEvent struct {
	Event   string `json:"event"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func NewEventBusNats() *EventBusNats {
	client, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}
	return &EventBusNats{client: client}
}

type EventBusNats struct {
	client *nats.Conn
}

func (e *EventBusNats) PublishEvent(kind string, user domain.User) error {
	u := UserEvent{
		Event:   kind,
		Id:      user.Id,
		Name:    user.Name,
		Surname: user.Surname,
	}
	content, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	return e.client.Publish("user", content)
}
