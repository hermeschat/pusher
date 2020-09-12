package nats

import (
	"github.com/hermeschat/proto"
	"github.com/nats-io/stan.go"
)

type MessageHandler interface {
	StreamSubscriber(subject string, g proto.Pusher_PusherEventBuffServer) (stan.Subscription, error)
	StreamPublisher(subject string, msg []byte) error
}

type tempStructNameNats struct {
	stan.Conn
}

func InitMessageHandler(clusterId string, clientId string) (MessageHandler, error) {
	c, err := StanStreamConnection(clusterId, clientId)
	if err != nil {
		return nil, err
	}
	return &tempStructNameNats{c}, nil
}

func StanStreamConnection(clusterId string, clientId string) (stan.Conn, error) {
	return stan.Connect(clusterId, clientId)
}

func (c *tempStructNameNats) StreamSubscriber(subject string, g proto.Pusher_PusherEventBuffServer) (stan.Subscription, error) {
	sub, err := c.Subscribe(subject, func(msg *stan.Msg) {
		_ = g.Send(&proto.PusherEvent{Event: &proto.PusherEvent_Server{Server: &proto.Server{
			FromSub: msg.Subject,
			Data:    msg.Data,
		}}})
		// TODO error need to be handled in a reasonable way
		// maybe cache it or send it to a worker for handling it
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *tempStructNameNats) StreamPublisher(subject string, msg []byte) error {
	return c.Publish(subject, msg)
}

func StreamUnsubscribe(sub stan.Subscription) error {
	return sub.Unsubscribe()
}
