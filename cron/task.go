package cron

import (
	"encoding/json"
	"mq/conf"
	"mq/dao"
	"mq/grpc"
	"mq/model"

	"mq/queue"
	. "mq/service"
	"mq/utils"

	"github.com/robfig/cron"
)

//结果补偿任务
type ToAck struct{}

//消息投递MQ任务
type ToMQ struct{}

//极端情况下解锁分布式锁
type ToUnLock struct{}

func (this ToAck) Run() {

	pushs, err := SrvHandle.GetTimeOutAck()

	if err != nil {
		// utils.WriteLog(dao.DBCHandle, pushs, err)
		return
	}

	//回调业务方
	for _, push := range *pushs {

		//请求
		address := push.Notify
		orderSn := push.OrderSn

		if push.Way == HTTP {

			// method := "GET"
			method := "POST"
			params := &map[string]interface{}{
				"order_sn": orderSn,
			}

			body, err := utils.RequstCleint(method, address, params)
			if err != nil {
				callBack(PREACK, orderSn)
				utils.WriteLog(dao.DBCHandle, params, err)
				return
			}

			resp := &model.Ack{}
			err = json.Unmarshal(body, resp)
			if err != nil {
				callBack(PREACK, orderSn)
				utils.WriteLog(dao.DBCHandle, resp, err)
				return
			}
			//更新数据

			callBack(resp.Ack, orderSn)

		} else if push.Way == GRPC {

			ack := push.Ack
			resp, err := grpc.Ask(address, orderSn, ack)
			if err != nil {
				callBack(REMOVEACK, orderSn)
				utils.WriteLog(dao.DBCHandle, resp, err)
				return
			}

			//更新数据
			callBack(resp.Ack, orderSn)
		}

	}

}

//定时推送
func (this ToMQ) Run() {

	pushs, err := SrvHandle.GetPush()
	if err != nil {
		// utils.WriteLog(dao.DBCHandle, pushs, err)
		return
	}

	//入队列
	for _, push := range *pushs {
		queue.Sdl.TaskJar.Push(&push, OP_REMOVE, OP_SORT)
	}
	return
}

//解锁
func (this ToUnLock) Run() {

	err := SrvHandle.InitLock()
	if err != nil {
		utils.WriteLog(dao.DBCHandle, "分布式锁", err)
		return
	}

	err = SrvHandle.InitPushLock()
	if err != nil {
		utils.WriteLog(dao.DBCHandle, "推送数据锁", err)
		return
	}

	return
}

func CronJob(c *conf.Config) {

	cr := cron.New()

	//回调补偿
	notifySpec := c.Cron.NotifySpec
	cr.AddJob(notifySpec, ToAck{})

	//MQ补偿
	sendSpec := c.Cron.SendSpec
	cr.AddJob(sendSpec, ToMQ{})

	//解锁
	lockSpec := c.Cron.LockSpec
	cr.AddJob(lockSpec, ToUnLock{})

	//启动计划任务
	cr.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer cr.Stop()

	select {}

}
