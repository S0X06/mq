package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mq/model"
	"net"

	pb "mq/grpc/answer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50052"
)

type server struct{} //服务对象

// SayHello 实现服务的接口 在proto中定义的所有服务都是接口
func (s *server) Prove(ctx context.Context, r *pb.Ack) (*pb.Response, error) {
	// log.Fatalf("Prove")
	fmt.Println("Prove:%s\n", r.OrderSn)

	ack := &model.Ack{
		OrderSn: r.GetOrderSn(),
		Ack:     2,
	}
	data, err := json.Marshal(ack)
	if err != nil {
		return &pb.Response{Code: 1000, Message: "OK"}, nil
	}

	return &pb.Response{Code: 200, Message: "OK", Data: data}, nil
}

func main() {

	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("监听失败：%v\n", err)
	}
	s := grpc.NewServer() //起一个服务
	pb.RegisterAnswerServer(s, &server{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		fmt.Println("服务开启失败：%v\n", err)
	}

}
