package main

import (
	"context"
	"demo/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

//实现服务端方法
type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "hello " + in.Name}, nil
}

func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	words := []string{
		"你好",
		"hello",
		"こんにちは",
		"안녕하세요",
	}
	for _, word := range words {
		data := &pb.HelloResponse{Message: word + in.GetName()}
		err := stream.Send(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	reply := "hello"
	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&pb.HelloResponse{Message: "hello " + reply})
			}
			return err
		}

		reply += recv.GetName() + ","
	}
}

func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {

	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		stream.SetTrailer(trailer)
	}()

	//读取流中的header
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(400, "no metadata")
	}
	token, ok := md["token"]
	if ok {
		for _, t := range token {
			fmt.Println(t)
		}
	}

	//创建和发送header
	header := metadata.New(map[string]string{
		"location": "beijing",
	})
	stream.SendHeader(header)

	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Fatalf("recv error: %v", err)
		}
		reply := magic(recv.GetName())
		err = stream.Send(&pb.HelloResponse{Message: reply})
		if err != nil {
			return err
		}
	}
}

func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	s := grpc.NewServer()                  // 创建gRPC服务端
	pb.RegisterGreeterServer(s, &server{}) // 注册服务方法到服务端
	//启动服务
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
