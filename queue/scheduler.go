package queue

import (
	"mq/model"
)

var (
	buf         = (1 << 10)
	mapCallBack = make(map[string]CallBack)
	Sdl         = NewScheduler()
)

type Scheduler struct {
	TaskJar        *TaskJar        //任务
	Conf           chan model.Conf //配置
	TmpMapCallBack chan map[string]func(orderSn string)
}

func NewScheduler() *Scheduler {

	taskJar := &TaskJar{
		SendTask: make(chan *SendTask, buf),
	}

	scheduler := &Scheduler{
		TaskJar:        taskJar,
		Conf:           make(chan model.Conf, buf),
		TmpMapCallBack: make(chan map[string]func(orderSn string), buf),
	}

	go scheduler.work()

	return scheduler

}

//异步发送
func (this *Scheduler) work() {

	for {

		select {

		case sendTask := <-this.TaskJar.Pull():
			//队列发送
			handle := Handle()

			go handle.Publish(
				sendTask.Exchange,
				sendTask.RoutingKey,
				sendTask.OrderSn,
				sendTask.Data,
				mapCallBack[sendTask.OPRemoveCode],
				mapCallBack[sendTask.OPSortCode],
			)

		case conf := <-this.Conf:
			//队列配置更新
			queue := &QueueConf{
				AppId:        conf.AppId,
				QueueName:    conf.QueueName,
				RoutingKey:   conf.RoutingKey,
				ExchangeName: conf.ExchangeName,
				ExchangeType: conf.ExchangeType,
				Done:         false,
			}

			Handle().AppendChannelConf(queue)

		case tmpFn := <-this.TmpMapCallBack:

			for opCode, fn := range tmpFn {
				if _, ok := mapCallBack[opCode]; !ok {
					// 添加
					mapCallBack[opCode] = fn
				}
			}

		default:

		}
	}

}
