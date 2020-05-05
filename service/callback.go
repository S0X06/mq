package service

import (
	"mq/dao"
	"mq/utils"

	"gopkg.in/mgo.v2/bson"
)

//即时发送队列回调
func (s *Service) OpPushRemove(orderSn string) {

	go func(orderSn string) {

		err := s.dao.RemovePush(orderSn)
		if err != nil {
			utils.WriteLog(dao.DBCHandle, "OpPushRemove,order_sn:"+orderSn, err)
		}

	}(orderSn)

	return
}

//回调可发送更新 -- 成功
func (s *Service) OpPushUpdate(orderSn string) {

	go func(orderSn string) {

		where := bson.M{
			"order_sn": orderSn,
			"ack":      PREACK,
		}

		data := bson.M{
			"$set": bson.M{
				"lock": UNLOCK,
				"sort": ZEROSORT,
				"ack":  ACK,
			},
		}

		err := s.dao.Set(where, data)
		if err != nil {
			// utils.WriteLog(dao.DBCHandle, "OpPushUpdate,order_sn:"+orderSn, err)
		}

	}(orderSn)

	return

}

//失败排序 -- sort 越大排序越后
func (s *Service) OpPushSort(orderSn string) {

	go func(orderSn string) {

		where := bson.M{
			"order_sn": orderSn,
		}

		data := bson.M{
			"$set": bson.M{
				"lock": UNLOCK,
			},
			"$inc": bson.M{
				"sort": INCSORT,
			},
		}

		err := s.dao.Set(where, data)
		if err != nil {
			utils.WriteLog(dao.DBCHandle, "OpPushSortInc,order_sn:"+orderSn, err)
		}

	}(orderSn)

	return

}

//解锁
func (s *Service) OpPushUnLock(orderSn string) {

	go func(orderSn string) {

		where := bson.M{
			"order_sn": orderSn,
		}

		data := bson.M{
			"$set": bson.M{
				"lock": UNLOCK,
			},
		}

		err := s.dao.Set(where, data)
		if err != nil {
			utils.WriteLog(dao.DBCHandle, "OpPushUnLock,order_sn:"+orderSn, err)
		}

	}(orderSn)

	return

}

//定时器批量发送队列回调
// func (s *Service) OpFilterPushRemove(orderSn string, opt ...interface{}) {

// 	go func(orderSn string, ack int) {

// 		err := s.dao.RemoveFilter(orderSn, ack)
// 		if err != nil {
// 			utils.WriteLog(dao.DBCHandle, "OpFilterPushRemove-RemoveFilter,order_sn:"+orderSn, err)
// 			return
// 		}

// 		err = s.dao.RemovePush(orderSn)
// 		if err != nil {
// 			utils.WriteLog(dao.DBCHandle, "OpFilterPushRemove-RemovePush,order_sn:"+orderSn, err)
// 		}

// 		return

// 	}(orderSn, ACK)

// 	return
// }

// //确认业务方删除回调
// func (s *Service) OpNotifyRemove(orderSn string, opt ...interface{}) {

// 	go func(orderSn string, ack int) {

// 		err := s.dao.RemoveFilter(orderSn, ack)
// 		if err != nil {
// 			utils.WriteLog(dao.DBCHandle, "OpNotifyRemove,order_sn:"+orderSn, err)
// 			return
// 		}

// 		err = s.dao.RemovePush(orderSn)
// 		if err != nil {
// 			utils.WriteLog(dao.DBCHandle, "OpNotifyRemove,order_sn:"+orderSn, err)
// 		}

// 		return

// 	}(orderSn, PREACK)

// 	return

// }

// //确认业务方更新回调
// func (s *Service) OpNotifyUpdate(orderSn string, opt ...interface{}) {

// 	go func(orderSn string, ack int) {

// 		err := s.dao.RemoveFilter(orderSn, ack)
// 		if err != nil {
// 			utils.WriteLog(dao.DBCHandle, "OpNotifyUpdate,order_sn:"+orderSn, err)
// 			return
// 		}

// 		//可发送
// 		err = s.dao.UpdateAck(orderSn, PREACK, ACK)
// 		if err != nil {
// 			utils.WriteLog(dao.DBCHandle, "OpNotifyUpdate,order_sn:"+orderSn, err)
// 		}

// 		return

// 	}(orderSn, PREACK)

// 	return
// }

// //删除过滤
// func (s *Service) OpFilterRemove(orderSn string, opt ...interface{}) {

// 	go func(orderSn string, ack int) {

// 		err := s.dao.RemoveFilter(orderSn, ack)
// 		if err != nil {
// 			utils.WriteLog(dao.DBCHandle, "OpFilterRemove,order_sn:"+orderSn, err)
// 		}

// 		return

// 	}(orderSn, 1)

// 	return

// }
