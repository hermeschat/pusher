package nats

import (
	"fmt"
	"github.com/nats-io/stan.go"
)

type INats interface {
	streamSubscriber(subject string) (stan.Subscription, error)
	streamPublisher(subject string, msg []byte) error
}

type TempStructNameNats struct {
	stan.Conn
}

func NatsStreamConnection() (stan.Conn, error)  {
	return stan.Connect("", "")
}

func (c *TempStructNameNats) streamSubscriber(subject string) (stan.Subscription, error) {
	sub, err := c.Subscribe(subject, func(msg *stan.Msg) {
		// TODO send message to subject owner
		fmt.Println(msg.Data)
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *TempStructNameNats) streamPublisher(subject string, msg []byte) error {
	return c.Publish(subject, msg)
}

func streamUnsubscribe(sub stan.Subscription) error {
	return sub.Unsubscribe()
}
