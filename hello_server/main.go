package main

import (
	"context"
	"demo/pb"
	"google.golang.org/grpc"
	"log"
	"net"
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
