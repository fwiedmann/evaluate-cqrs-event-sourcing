package infra

import (
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user-query/domain"
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

func (e EventBusNats) Subscribe(handler func(e domain.Event) error) error {
	_, err := e.client.Subscribe("user", func(msg *nats.Msg) {
		fmt.Println("EVENT OCCURRED")
		u := &UserEvent{}
		if err := json.Unmarshal(msg.Data, u); err != nil {
			fmt.Println(err.Error())
			return
		}

		err := handler(domain.Event{
			Kind: domain.EventKind(u.Event),
			User: domain.User{
				Id:      u.Id,
				Name:    u.Name,
				Surname: u.Surname,
			},
		})
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
