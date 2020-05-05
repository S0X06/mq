package model

type (
	Ack struct {
		OrderSn string `form:"order_sn" json:"order_sn" binding:"required"`
		Ack     int32  `form:"ack" json:"ack" binding:"required"` //生产者 ： 0：删除 1:预发送,2:可发送
	}

	Receive struct {
		AppId   string `form:"app_id" json:"app_id" binding:"required"`
		Notify  string `form:"notify" binding:"required"`
		OrderSn string `form:"order_sn" json:"order_sn" binding:"required"`
		Data    string `form:"data" json:"data" binding:"required"`
		Way     int    `form:"way" `
	}
)
