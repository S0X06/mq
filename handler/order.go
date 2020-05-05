package handler

import (
	"mq/model"
	. "mq/service"
	"mq/utils"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
)

//预发送 => 状态待发送
func Try(c *gin.Context) {

	rec := &model.Receive{}

	if err := c.ShouldBind(rec); err != nil {
		// fmt.Println(err.Error())
		SendResponse(utils.FAILURE, c, err.Error())
		return

	}
	// fmt.Println(rec)

	//设置网络请求
	rec.Way = HTTP
	err := SrvHandle.Try(rec)

	if err != nil {
		SendResponse(utils.FAILURE, c, err)
		return

	}

	SendResponse(utils.SUCCESS, c, "")
	return

}

//是否可发送
func PublisherAck(c *gin.Context) {

	ack := &model.Ack{}

	if err := c.ShouldBind(ack); err != nil {
		SendResponse(utils.FAILURE, c, err.Error())
		return

	}

	if ack.Ack == ACK {

		// fmt.Println(ack)
		err := SrvHandle.Ack(ack)

		if err != nil {
			// fmt.Println(err)
			SendResponse(utils.FAILURE, c, err)
			return

		}

	} else if ack.Ack == REMOVEACK {

		err := SrvHandle.RemoveAck(ack)

		if err != nil {
			// fmt.Println(err)
			SendResponse(utils.FAILURE, c, err)
			return
		}

	} else {

		SendResponse(utils.FAILURE, c, "ack 错误")
		return
	}

	SendResponse(utils.SUCCESS, c, "成功")
	return

}
