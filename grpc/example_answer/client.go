package main

import (
	"context"
	"log"
	"time"

	pb "mq/grpc/answer"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	//建立链接

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAnswerClient(conn)

	// Contact the server and print out its response.
	// name := defaultName
	// if len(os.Args) > 1 {
	// 	name = os.Args[1]
	// }
	// 1秒的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Prove(ctx, &pb.Ack{OrderSn: defaultName, Ack: 2})
	if err != nil {
		log.Fatalf("could not reply: %v", err)
	}
	log.Printf("reply: %s", r.Message)
}
