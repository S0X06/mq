package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"mq/conf"
	"time"

	"sync"

	"mq/dao"
	"mq/utils"

	"github.com/streadway/amqp"
)

const (
	HEART = (1 << 1)
)

var (
	Lock       sync.RWMutex
	amqpHandle = &AmqpChannel{}
)

type (
	CallBack func(orderSn string)

	QueueConf struct {
		AppId        string
		QueueName    string
		RoutingKey   string
		ExchangeName string
		ExchangeType string
		Done         bool
	}

	AmqpChannel struct {
		AmqpAddr        string
		TmpDone         bool
		Done            chan bool
		Conn            *amqp.Connection
		Channel         *amqp.Channel
		NotifyConnClose chan *amqp.Error
		NotifyChanClose chan *amqp.Error
		QueueConf       *QueueConf
		DoneChannel     chan bool
	}
)

//逻辑通道

func Init(conf *conf.Config) (err error) {

	rabbitConf := conf.RabbitMq

	amqpHandle, err = NewAmqpManager(rabbitConf)

	return err

}

//创建mq通道
func NewAmqpManager(conf *conf.RabbitMq) (*AmqpChannel, error) {

	amqpAddr := fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.UserName, conf.PassWord, conf.Addr, conf.Port)
	fmt.Println(amqpAddr)
	manager := &AmqpChannel{
		AmqpAddr:        amqpAddr,
		TmpDone:         false,
		Done:            make(chan bool, 2),
		NotifyConnClose: make(chan *amqp.Error),
		NotifyChanClose: make(chan *amqp.Error),
		DoneChannel:     make(chan bool, 2),
	}

	go manager.connect()

	return manager, nil
}

func Handle() *AmqpChannel {
	return amqpHandle
}

func (this *AmqpChannel) connect() {

	for {

		select {

		case done := <-this.Done:

			if !done {

				conn, err := amqp.Dial(this.AmqpAddr)
				if err != nil {
					fmt.Printf("MQ 连接失败,失败原因： %s \n, 2 秒后尝试重新连接....\n", HEART, err.Error())
					continue
				}

				conn.NotifyClose(this.NotifyConnClose)

				//通道
				ch, err := conn.Channel()
				if err != nil {
					fmt.Printf("MQ 通道打开失败,失败原因： %s \n, 2 秒后尝试重新连接....\n", HEART, err.Error())
					continue
				}

				ch.NotifyClose(this.NotifyChanClose)

				this.Channel = ch
				this.Conn = conn
				this.TmpDone = true
				this.Done <- true

			} else {
				//关闭 true  否则false
				if !this.Conn.IsClosed() {
					fmt.Printf("MQ 连接成功 ...\n")
				}

			}

		case connClose := <-this.NotifyConnClose:
			fmt.Printf("NotifyConnClose err:%v\n", connClose)
			this.TmpDone = false
			this.Done <- false

		case chanClose := <-this.NotifyChanClose:
			fmt.Printf("NotifyChanClose err:%v\n", chanClose)
			this.TmpDone = false

		case doneChannel := <-this.DoneChannel:

			//声明通道
			if done := this.TmpDone; done && doneChannel { //未连接,发送会出现错误
				this.Declare()
			}

		default:

			time.Sleep(HEART * time.Second)
			if !this.TmpDone || this.Conn == nil || this.Channel == nil {
				this.TmpDone = false
				this.Done <- false
			}

		}

	}

}

//关闭连接
func (this *AmqpChannel) Close() {
	this.Channel.Close()
	this.Conn.Close()
}

//声明队列/路由/交换机
func (this *AmqpChannel) Declare() {

	if !this.QueueConf.Done {
		//声明交换机
		err := this.Channel.ExchangeDeclare(
			this.QueueConf.ExchangeName,
			this.QueueConf.ExchangeType, //Direct:精确(一对一) , Topic:(模糊一对多), Fanout(无路由键,一对多广播所有关联队列), Headers(无路由键,一对多发送队列，通过header键值匹配到才发送)
			true,                        //持久化
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			return
		}

		//声明队列
		// _, err = this.Channel.QueueDeclare(
		// 	this.QueueConf.QueueName,
		// 	true,
		// 	false,
		// 	false,
		// 	false,
		// 	nil,
		// )
		// if err != nil {
		// 	return
		// }

		// //绑定队列
		// err = this.Channel.QueueBind(
		// 	this.QueueConf.QueueName,
		// 	this.QueueConf.RoutingKey,
		// 	this.QueueConf.ExchangeName,
		// 	false,
		// 	nil,
		// )
		// if err != nil {
		// 	return
		// }

		this.QueueConf.Done = true
	}

	return
}

//声明队列/路由/交换机
func (this *AmqpChannel) AppendChannelConf(conf *QueueConf) {

	this.QueueConf = conf
	this.DoneChannel <- true

}

//推送
func (this *AmqpChannel) Publish(exchange, routingkey, orderSn string, data string, removeCallBack CallBack, sortCallBack CallBack) error {
	// 判断channel是否正常

	if done := this.TmpDone; !done {
		sortCallBack(orderSn)
		return errors.New("通道不可用")
	}

	body, err := json.Marshal(data)
	if err != nil {
		sortCallBack(orderSn)

		utils.WriteLog(dao.DBCHandle, orderSn, err)
		return err
	}

	content := amqp.Publishing{
		DeliveryMode: amqp.Persistent, //消息持久化(不是完全强一致性),基于队列持久化
		ContentType:  "text/plain",
		Body:         body,
	}

	//开启事物
	// err = this.Channel.Tx()
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println("cccc")
	// 	sortCallBack(orderSn)
	// 	utils.WriteLog(dao.DBCHandle, orderSn, err)
	// 	return err
	// }

	//发布
	err = this.Channel.Publish(exchange, routingkey, false, false, content)
	if err != nil {
		//回滚
		// this.Channel.TxRollback()

		sortCallBack(orderSn)
		utils.WriteLog(dao.DBCHandle, orderSn, err)
		return err
	}

	// //提交
	// err = this.Channel.TxCommit()
	// if err != nil {
	// 	sortCallBack(orderSn)
	// 	utils.WriteLog(dao.DBCHandle, orderSn, err)
	// 	return err
	// }

	removeCallBack(orderSn)

	return nil

}
