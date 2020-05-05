package handler

import (
	"mq/model"
	. "mq/service"

	"github.com/gin-gonic/gin"
)

func AddConf(c *gin.Context) {

	conf := &model.Conf{}

	if err := c.ShouldBind(conf); err != nil {
		SendResponse(-1, c, err.Error())
		return
	}

	err := SrvHandle.CreateConf(conf)

	if err != nil {
		SendResponse(-1, c, err)
		return
	}

	SendResponse(0, c, "")
	return

}

//更新
func UpdateConf(c *gin.Context) {

	conf := &model.Conf{}

	if err := c.ShouldBind(conf); err != nil {
		SendResponse(-1, c, err.Error())
		return
	}

	err := SrvHandle.UpdateConf(conf)

	if err != nil {
		SendResponse(-1, c, err)
		return
	}

	SendResponse(0, c, "")
	return

}

//确认发布
func ReleaseConf(c *gin.Context) {

	conf := &model.Conf{}

	if err := c.ShouldBind(conf); err != nil {
		SendResponse(-1, c, err.Error())
		return
	}

	conf.Status = 1
	err := SrvHandle.UpdateConf(conf)

	if err != nil {
		SendResponse(-1, c, err)
		return
	}

	SendResponse(0, c, "")
	return

}

//获取所有配置
func GetConf(c *gin.Context) {

	date, err := SrvHandle.GetAllConf()

	if err != nil {
		SendResponse(-1, c, err)
		return
	}

	SendResponse(0, c, date)
	return

}

//删除
func RemoveConf(c *gin.Context) {

	conf := &model.Conf{}

	if err := c.ShouldBind(conf); err != nil {
		SendResponse(-1, c, err.Error())
		return
	}

	err := SrvHandle.RemoveConf(conf.AppId)

	if err != nil {
		SendResponse(-1, c, err)
		return
	}

	SendResponse(0, c, "")
	return

}
