package main

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hello_client/pb"
	"log"
)

var (
	addr = flag.String("addr", ":8080", "address to listen")
	name = flag.String("name", "leon", "name to say hello to")
)

func main() {
	flag.Parse()
	//连接的server端，禁用安全传输
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	//调用服务端方法unary
	//ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	//defer cancelFunc()
	//resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
	//if err != nil {
	//	log.Fatalf("say hello error: %v", err)
	//}
	//fmt.Print(resp.Message)

	//调用服务端流式方法
	//runLotsOfReplies(client)

	//runLotsOfGreeting(client)

	runBidiHello(client)
}
