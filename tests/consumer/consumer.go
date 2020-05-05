package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"mq/conf"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

const (
	HEART = (1 << 1)
)

var (
	userName = ""
	passWord = ""
	addr     = ""
	port     = ""

	ExchangeName = "demo"
	ExchangeType = "topic"
	QueueName    = "demo"
	RoutingKey   = "demo"
)

type AmqpChannel struct {
	AmqpAddr        string
	TmpDone         bool
	Done            chan bool
	Conn            *amqp.Connection
	Channel         *amqp.Channel
	NotifyConnClose chan *amqp.Error
	NotifyChanClose chan *amqp.Error
}

//逻辑通道

func main() {

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	userName = conf.Conf.RabbitMq.UserName
	passWord = conf.Conf.RabbitMq.PassWord
	addr = conf.Conf.RabbitMq.Addr
	port = conf.Conf.RabbitMq.Port

	_, _ = NewAmqpManager()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		fmt.Printf("get a signal : %s \n", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

//创建mq通道
func NewAmqpManager() (*AmqpChannel, error) {

	amqpAddr := fmt.Sprintf("amqp://%s:%s@%s:%s/", userName, passWord, addr, port)
	fmt.Println(amqpAddr)
	manager := &AmqpChannel{
		AmqpAddr:        amqpAddr,
		TmpDone:         false,
		Done:            make(chan bool, 2),
		NotifyConnClose: make(chan *amqp.Error),
		NotifyChanClose: make(chan *amqp.Error),
	}

	go manager.connect()

	return manager, nil
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

				fmt.Printf("MQ 连接成功...\n")

			} else {

				this.Declare()
			}

		case connClose := <-this.NotifyConnClose:
			fmt.Printf("NotifyConnClose err:%v\n", connClose)
			this.TmpDone = false
			this.Done <- false

		case chanClose := <-this.NotifyChanClose:
			fmt.Printf("NotifyChanClose err:%v\n", chanClose)
			this.TmpDone = false

		default:

			if !this.TmpDone || this.Conn == nil {
				time.Sleep(HEART * time.Second)
				this.TmpDone = false
				this.Done <- false
			} else {
				this.Consume()
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

	if this.TmpDone {
		//声明交换机
		err := this.Channel.ExchangeDeclare(
			ExchangeName,
			ExchangeType,
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			return
		}

		//声明队列
		_, err = this.Channel.QueueDeclare(
			QueueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return
		}

		//绑定队列
		err = this.Channel.QueueBind(
			QueueName,
			RoutingKey,
			ExchangeName,
			false,
			nil,
		)
		if err != nil {
			return
		}

		err = this.Channel.Qos(
			1,
			0,
			false,
		)

		if err != nil {
			return
		}

	}

	return
}

//消费
func (this *AmqpChannel) Consume() {

	amqpBody, err := this.Channel.Consume(
		QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
	}

	for {

		select {
		case body := <-amqpBody:
			var resp string
			err = json.Unmarshal(body.Body, &resp)
			if err != nil {
				fmt.Println("err:", err)
				body.Nack(false, true)
				continue
			}
			//已消费
			body.Ack(true)
			fmt.Println("body:", resp)
		}
	}

	// return nil
}
