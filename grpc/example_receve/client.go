package main

import (
	"context"
	"fmt"
	"log"
	"mq/utils"
	"time"

	pb "mq/grpc/receive"

	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
	notify  = "127.0.0.1:50052"
	appId   = "123"
)

func send(conn *grpc.ClientConn) {

	c := pb.NewRecipientClient(conn)
	// 1秒的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	receive := &pb.Receive{
		OrderSn: utils.GenValidateCode(10),
		AppId:   appId,
		Notify:  notify,
		Data:    utils.GenValidateCode(5),
	}

	fmt.Printf("send:%s\n", receive.OrderSn)
	r, err := c.Try(ctx, receive)
	if err != nil {
		fmt.Println("预发送错误: %v", err)
	}

	if r.Code == 0 {
		okAck := &pb.Ack{
			OrderSn: receive.OrderSn,
			Ack:     2,
		}
		_, err := c.PublisherAck(ctx, okAck)
		if err != nil {
			fmt.Println("可发送错误: %v", err)
		}

		return
	}
	// fmt.Println("删除错误: %v", r.Message)

	ack := &pb.Ack{
		OrderSn: receive.OrderSn,
		Ack:     0,
	}
	_, err = c.PublisherAck(ctx, ack)

	if err != nil {
		fmt.Println("删除错误: %v", err)
	}
}

func main() {
	//建立链接

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	defer conn.Close()

	for {
		go send(conn)
		time.Sleep(time.Second)
	}

	do := make(chan os.Signal, 1)
	signal.Notify(do, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-do

		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:

			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}
