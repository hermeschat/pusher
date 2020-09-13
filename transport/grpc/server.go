package grpc

import (
	"fmt"
	"github.com/hermeschat/proto"
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
	mh nats.MessageHandler
}

func (ps *PusherService) PusherEventBuff(server proto.Pusher_PusherEventBuffServer) error {
	// TODO client id must change into user id or sth from user identifier
	fmt.Println("new user connected")
	nconn, err := nats.InitMessageHandler("test-cluster", "change_this")
	if err != nil {
		return err
	}

	ps.mh = nconn
	for {
		event, err := server.Recv()
		if err != nil {
			continue
		}
		if _, msg := event.Event.(*proto.PusherEvent_Client); msg {
			err = ps.mh.StreamPublisher(event.GetClient().ToSub, event.GetClient().Data)
			if err != nil {
				// handle errors in here in better way
			}
			fmt.Println("published")
			continue
		} else if _, msg := event.Event.(*proto.PusherEvent_Join); msg {
			// TODO handle subscriber variable maybe store address or cache it
			_, err := ps.mh.StreamSubscriber(event.GetJoin().MySub, server)
			if err != nil {
				// TODO handle error
			}
			fmt.Println("joined")
		} else {
			// TODO send not supported event
		}
	}
}
