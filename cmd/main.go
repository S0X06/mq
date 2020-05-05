package main

import (
	"flag"
	"mq/conf"

	"mq/cron"

	"mq/queue"
	"mq/router"
	"mq/service"
	"os"
	"os/signal"
	"syscall"

	"mq/grpc"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	logrus.Info(" start...")

	//注册etcd

	//数据库连接
	err := service.Init(conf.Conf)
	if err != nil {
		logrus.Info("数据库连接失败...,失败原因：", err)
		// panic(err)
	}

	//消息队列连接
	err = queue.Init(conf.Conf)
	if err != nil {
		logrus.Info("队列连接失败...,失败原因：", err)
	}

	//开启 grpc 接口
	go grpc.Server(conf.Conf.Grpc.Port)

	//定时器
	go cron.CronJob(conf.Conf)

	//路由
	engine := router.InitRouter()

	//设置模式
	gin.SetMode(conf.Conf.RunMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	engine.Run(conf.Conf.Addr)

	//平滑重启
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		logrus.Info("get a signal : %s ", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logrus.Info("exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}
