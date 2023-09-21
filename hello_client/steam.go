package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"hello_client/pb"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func runLotsOfReplies(c pb.GreeterClient) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	stream, err := c.LotsOfReplies(ctx, &pb.HelloRequest{Name: "leon"})
	if err != nil {
		log.Fatalf("say hello error: %v", err)
	}
	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("recv error: %v", err)
			}
		}
		fmt.Printf("%s\n", recv.Message)
	}
}

func runLotsOfGreeting(c pb.GreeterClient) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		log.Fatalf("say hello error: %v", err)
	}
	names := []string{
		"leon", "liqixuan", "lqx",
	}
	for _, name := range names {
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("send error: %v", err)
		}
	}
	recv, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("recv error: %v", err)
	}
	fmt.Printf("%s\n", recv.Message)
}

func runBidiHello(c pb.GreeterClient) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFunc()

	md := metadata.Pairs("token", "123456")
	ctx = metadata.NewOutgoingContext(ctx, md)
	
	stream, err := c.BidiHello(ctx)
	if err != nil {
		log.Fatalf("say hello error: %v", err)
	}
	waitC := make(chan struct{})
	go func() {
		for {
			recv, err2 := stream.Recv()
			if err2 != nil {
				if err2 == io.EOF {
					close(waitC)
					return
				} else {
					log.Fatalf("recv error: %v", err2)
				}
			}
			fmt.Println(recv.Message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("read string error: %v", err)
		}
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		if cmd == "exit" {
			break
		}
		err = stream.Send(&pb.HelloRequest{Name: cmd})
		if err != nil {
			log.Fatalf("send error: %v", err)
		}
	}
	stream.CloseSend()
	<-waitC
}
