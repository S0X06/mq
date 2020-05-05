package grpc

import (
	"context"
	"errors"
	"mq/dao"
	"mq/model"
	. "mq/service"
	"mq/utils"
	"net"

	pb "mq/grpc/receive"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{} //服务对象

//  预发送
func (s *server) Try(ctx context.Context, r *pb.Receive) (*pb.Response, error) {

	rec := &model.Receive{
		AppId:   r.GetAppId(),
		Way:     GRPC,
		Notify:  r.GetNotify(),
		OrderSn: r.GetOrderSn(),
		Data:    r.GetData(),
	}

	err := SrvHandle.Try(rec)

	if err != nil {
		return &pb.Response{Code: utils.FAILURE}, err
	}

	return &pb.Response{Code: utils.SUCCESS}, nil
}

//是否可推送
func (s *server) PublisherAck(ctx context.Context, r *pb.Ack) (*pb.Response, error) {

	ack := &model.Ack{
		OrderSn: r.GetOrderSn(),
		Ack:     r.GetAck(),
	}

	if ack.Ack == ACK {

		err := SrvHandle.Ack(ack)

		if err != nil {
			return &pb.Response{Code: utils.FAILURE}, err
		}

	} else if ack.Ack == REMOVEACK {

		err := SrvHandle.RemoveAck(ack)

		if err != nil {
			return &pb.Response{Code: utils.FAILURE}, err
		}

	} else {

		return &pb.Response{Code: utils.FAILURE}, errors.New("ack错误")
	}

	return &pb.Response{Code: utils.SUCCESS}, nil
}

//服务
func Server(port string) {

	listen, err := net.Listen("tcp", port)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, "GRPC 监听端口失败", err)
		panic(err)
	}
	s := grpc.NewServer() //起一个服务
	pb.RegisterRecipientServer(s, &server{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		utils.WriteLog(dao.DBCHandle, "GRPC 注册失败", err)
		panic(err)
	}

}
