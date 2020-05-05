package model

import (
	"time"
)

type Push struct {
	OrderSn      string `bson:"order_sn" binding:"required"`
	AppId        string `bson:"app_id" binding:"required"`
	Ack          int32  `bson:"ack" binding:"required"` // GRPC 1  HTTP 2
	Notify       string `bson:"notify" binding:"required"`
	Data         string `bson:"data" binding:"required"`
	RoutingKey   string `bson:"routing_key"`
	ExchangeName string `bson:"exchange_name"`
	Way          int    `bson:"way"`
	Sort         int    `bson:"sort"`
	Lock         int    `bson:"lock"`
	LockAt       int64  `bson:"lock_at"`
	UpdatedAt    int64  `bson:"updated_at"`
	CreatedAt    int64  `bson:"created_at"`
}

//解析push filter
func ParsePushToFilter(pushs *[]Push, ack int) ([]interface{}, bool) {

	var filters []interface{} = make([]interface{}, 0)

	for _, push := range *pushs {

		filter := &Filter{
			OrderSn:   push.OrderSn,
			Ack:       ack,
			CreatedAt: time.Now().Unix(),
		}

		filters = append(filters, filter)

	}

	if len(filters) == 0 {
		return filters, false
	}

	return filters, true
}

//上锁
func ParsePush(pushs *[]Push) ([]string, bool) {

	var filters []string

	for _, push := range *pushs {

		filters = append(filters, push.OrderSn)

	}

	if len(filters) == 0 {
		return filters, false
	}

	return filters, true
}
