package queue

import (
	"mq/model"
)

type (
	SendTask struct {
		OrderSn      string `json:"order_sn"`
		Exchange     string `json:"exchange"`
		RoutingKey   string `json:"routing_key"`
		Data         string `json:"data"`
		OPRemoveCode string `json:"op_remove_code"`
		OPSortCode   string `json:"op_sort_code"`
	}

	TaskJar struct {
		SendTask chan *SendTask //发送数据
	}
)

//添加一个任务
func (this *TaskJar) Push(push *model.Push, OPRemoveCode string, OPSortCode string) {

	sendTask := &SendTask{
		OrderSn:      push.OrderSn,
		Data:         push.Data,
		RoutingKey:   push.RoutingKey,
		Exchange:     push.ExchangeName,
		OPRemoveCode: OPRemoveCode,
		OPSortCode:   OPSortCode,
	}

	this.SendTask <- sendTask
}

//获取一个任务
func (this *TaskJar) Pull() chan *SendTask {
	return this.SendTask
}
