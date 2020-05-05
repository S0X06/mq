package grpc

import (
	"context"
	"encoding/json"
	"mq/model"
	"time"

	pb "mq/grpc/answer"

	"google.golang.org/grpc"
)

func Ask(address string, orderSn string, ack int32) (*model.Ack, error) {

	resp := &model.Ack{}
	//建立链接
	// grpc.WithBalancerName() etcd
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return resp, err
	}
	defer conn.Close()
	c := pb.NewAnswerClient(conn)

	// 1秒的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	params := &pb.Ack{OrderSn: orderSn, Ack: ack}

	r, err := c.Prove(ctx, params)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(r.Data, resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
