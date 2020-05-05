package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "mq/grpc/receive"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{} //服务对象

// SayHello 实现服务的接口 在proto中定义的所有服务都是接口
func (s *server) Try(ctx context.Context, r *pb.Receive) (*pb.Response, error) {
	fmt.Println("Try")
	return &pb.Response{Code: 200, Message: "OK"}, nil
}

func (s *server) PublisherAck(ctx context.Context, r *pb.Ack) (*pb.Response, error) {
	fmt.Println("PublisherAck")
	return &pb.Response{Code: 200, Message: "OK"}, nil
}

func (s *server) Remove(ctx context.Context, r *pb.Ack) (*pb.Response, error) {
	fmt.Println("Remove")
	return &pb.Response{Code: 200, Message: "OK"}, nil
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //起一个服务
	pb.RegisterRecipientServer(s, &server{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
