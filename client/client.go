package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/hermeschat/proto"
	"google.golang.org/grpc"
	"os"
	"strings"
	"sync"
)

func main() {
	con, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "9000"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli:= proto.NewPusherClient(con)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println("client created")
	buff, err := cli.PusherEventBuff(ctx)
	if err != nil {
		panic(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your ID: ")
	mysub, _ := reader.ReadString('\n')
	mysub = mysub[:len(mysub) - 1]
	err = buff.Send(&proto.PusherEvent{Event: &proto.PusherEvent_Join{Join: &proto.Join{MySub: mysub}}})
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			msg, _ := reader.ReadString('\n')
			if msg == "" || msg == "\n" {
				continue
			}
			to := strings.SplitN(msg, " ", 2)[0]
			msg = strings.SplitN(msg, " ", 2)[1]
			msg = msg[: len(msg) - 1]
			msg = mysub + "::" + msg
			err = buff.Send(&proto.PusherEvent{Event: &proto.PusherEvent_Client{Client: &proto.Client{ToSub: to, Data: []byte(msg)}}})
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		for {
			pEvent, err := buff.Recv()
			if err != nil {
				panic(err.Error())
			}
			switch pEvent.GetEvent().(type) {
			case *proto.PusherEvent_Server:
				s := pEvent.GetServer()
				fmt.Printf("\n%s : %s", s.FromSub, strings.Split(string(s.Data), "::")[1])
			}
		}
	}()
	wg.Wait()
}