package infra

import (
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/core"
	"github.com/nats-io/nats.go"
)

type UserEvent struct {
	Event   string `json:"event"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type EventHandler interface {
	OnUserCreated(event core.UserCreatedEvent) error
	OnUserUpdated(event core.UserUpdatedEvent) error
	OnUserDeleted(event core.UserDeletedEvent) error
}

func NewEventBusNats(handler EventHandler) *EventBusNats {
	client, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}
	return &EventBusNats{client: client, handler: handler}
}

type EventBusNats struct {
	client  *nats.Conn
	handler EventHandler
}

func (e *EventBusNats) Subscribe() error {
	_, err := e.client.Subscribe("user", func(msg *nats.Msg) {
		fmt.Println("EVENT OCCURRED")
		u := &UserEvent{}
		if err := json.Unmarshal(msg.Data, u); err != nil {
			fmt.Println(err.Error())
			return
		}
		var err error
		switch u.Event {
		case core.UserCreatedEventKey:
			err = e.handler.OnUserCreated(core.UserCreatedEvent{User: mapUserEventToUser(*u)})
		case core.UserUpdatedEventKey:
			err = e.handler.OnUserUpdated(core.UserUpdatedEvent{User: mapUserEventToUser(*u)})
		case core.UserDeletedEventKey:
			err = e.handler.OnUserDeleted(core.UserDeletedEvent{UserId: u.Id})
		}
		if err != nil {
			msg.Nak()
			return
		}
		msg.Ack()
	},
	)
	fmt.Println("subscribed to user")
	return err
}

func mapUserEventToUser(event UserEvent) core.User {
	return core.User{
		Id:      event.Id,
		Name:    event.Name,
		Surname: event.Surname,
	}
}
