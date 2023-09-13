package main

import (
	"context"
	"fmt"
	"hello_client/pb"
	"io"
	"log"
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
