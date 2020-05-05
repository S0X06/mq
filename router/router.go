package router

import (
	"mq/handler"
	// "user/router/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	engine := gin.Default()
	engine.POST("/try", handler.Try)
	engine.PUT("/ack", handler.PublisherAck)

	engine.GET("/conf", handler.GetConf)
	engine.POST("/conf", handler.AddConf)
	engine.PUT("/conf", handler.UpdateConf)
	engine.DELETE("/conf", handler.RemoveConf)
	engine.PUT("/conf/release", handler.ReleaseConf) //发布队列

	return engine
}
