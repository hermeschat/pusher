package grpc

import (
	"fmt"
	"github.com/hermeschat/proto"
	"github.com/nats-io/stan.go"
	"google.golang.org/grpc"
	"net"
	"pusher/transport/nats"
)

func CreateGRPCServer() {
	listen, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	ps := PusherService{}
	proto.RegisterPusherService(server, &proto.PusherService{PusherEventBuff: ps.PusherEventBuff})
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}

type PusherService struct {
	proto.UnstablePusherService
}

type ClientProperty struct {
	mh nats.MessageHandler
	ssub stan.Subscription
	subject string
}

func (ps *PusherService) PusherEventBuff(server proto.Pusher_PusherEventBuffServer) error {
	// TODO client id must change into user id or sth from user identifier
	fmt.Println("new user connected")

	cp := new(ClientProperty)
	for {
		event, err := server.Recv()
		if err != nil {
			continue
		}
		if _, msg := event.Event.(*proto.PusherEvent_Client); msg {
			if cp.mh.IsConnected() != true {
				err := server.SendMsg("you have to join first")
				if err != nil {

				}
				continue
			}
			myBadge := []byte(cp.subject + "::")
			err = cp.mh.StreamPublisher(event.GetClient().ToSub, append(myBadge, event.GetClient().Data...))
			if err != nil {
				// handle errors in here in better way
			}
			fmt.Println("published")
			continue
		} else if _, msg := event.Event.(*proto.PusherEvent_Join); msg {
			nconn, err := nats.InitMessageHandler("test-cluster", event.GetJoin().MySub)
			if err != nil {
				return err
			}
			cp.mh = nconn
			// TODO handle subscriber variable maybe store address or cache it
			mysub := event.GetJoin().MySub

			sub, err := cp.mh.StreamSubscriber(mysub, server)
			if err != nil {
				// TODO handle error
			}
			cp.ssub = sub
			cp.subject = mysub
			fmt.Println("joined")
		} else {
			// TODO send not supported event
		}
	}
}
