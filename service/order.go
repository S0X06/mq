package service

import (
	"errors"
	"mq/dao"
	"mq/model"
	. "mq/queue"
	"mq/utils"
	"time"

	"gopkg.in/mgo.v2"
)

//生产者预确认
func (s *Service) Try(rec *model.Receive) error {

	conf, err := s.GetConf(rec.AppId)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, conf, err)
		return err
	}

	//获取数据
	_, err = s.dao.GetDataBySn(rec.OrderSn)
	if err != mgo.ErrNotFound {
		utils.WriteLog(dao.DBCHandle, rec, ErrFound)
		return ErrFound
	}

	t := time.Now().Unix()
	p := &model.Push{
		AppId:        rec.AppId,
		Notify:       rec.Notify,
		OrderSn:      rec.OrderSn,
		Data:         rec.Data,
		Ack:          PREACK,
		RoutingKey:   conf.RoutingKey,
		ExchangeName: conf.ExchangeName,
		Way:          rec.Way,
		Sort:         0,
		Lock:         0,
		UpdatedAt:    t,
		CreatedAt:    t,
	}
	err = s.dao.CreatePush(p)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, p, err)
	}
	return err

}

//生产者确认 ack
func (s *Service) Ack(ack *model.Ack) (err error) {

	// Ack == 2
	push, err := s.dao.GetDataBySn(ack.OrderSn)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, ack, err)
		return
	}
	//可发送
	err = s.dao.UpdateAck(ack.OrderSn, PREACK, ACK)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, ack, err)
		return
	}

	//入队列
	go Sdl.TaskJar.Push(push, OP_REMOVE, OP_SORT)

	return

}

//生产者确认 ack
func (s *Service) RemoveAck(ack *model.Ack) (err error) {
	//成功删除
	err = s.dao.RemovePush(ack.OrderSn)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, ack, err)
	}
	return
}

//获取推送数据
func (s *Service) GetPush() (*[]model.Push, error) {

	Lock.Lock()
	defer Lock.Unlock()

	//多程序抢锁
	err := ProcessLock(SENDLOCK, UNLOCK, LOCK)
	if err != nil {
		return &[]model.Push{}, err
	}

	//获取数据
	// var err error
	SendPushs, err = s.dao.GetPush(UNLOCK, ACK, LIMIT)
	if err != nil {
		ProcessLock(SENDLOCK, LOCK, UNLOCK)

		utils.WriteLog(dao.DBCHandle, SendPushs, err)
		return &[]model.Push{}, err
	}
	//上锁
	lockItems, ok := model.ParsePush(SendPushs)
	if !ok {
		ProcessLock(SENDLOCK, LOCK, UNLOCK)

		return &[]model.Push{}, errors.New("无推送数据")
	}

	err = s.dao.BatchLock(lockItems, LOCK)
	if err != nil {
		ProcessLock(SENDLOCK, LOCK, UNLOCK)

		return &[]model.Push{}, err
	}

	err = ProcessLock(SENDLOCK, LOCK, UNLOCK)
	if err != nil {
		return &[]model.Push{}, err
	}

	return SendPushs, err
}

//需要和业务确认的数据
func (s *Service) GetTimeOutAck() (*[]model.Push, error) {

	Lock.Lock()
	defer Lock.Unlock()

	//多程序抢锁
	err := ProcessLock(NOTIFYLOCK, UNLOCK, LOCK)
	if err != nil {
		return &[]model.Push{}, err
	}

	//获取数据
	// var err error
	RequstPushs, err = s.dao.GetTimeOutAck(PREACK, UNLOCK, LIMIT, TIMEOUT)
	if err != nil {
		ProcessLock(NOTIFYLOCK, LOCK, UNLOCK)
		// utils.WriteLog(dao.DBCHandle, pushs, err)
		return &[]model.Push{}, err
	}
	//上锁
	lockItems, ok := model.ParsePush(RequstPushs)
	if !ok {
		s.dao.ProcessLock(NOTIFYLOCK, LOCK, UNLOCK)

		return &[]model.Push{}, errors.New("无回调数据")
	}

	err = s.dao.BatchLock(lockItems, LOCK)
	if err != nil {

		ProcessLock(NOTIFYLOCK, LOCK, UNLOCK)

		return &[]model.Push{}, err
	}

	err = ProcessLock(NOTIFYLOCK, LOCK, UNLOCK)
	if err != nil {
		return &[]model.Push{}, err
	}

	return RequstPushs, err
}

//通过id获取消息
func (s *Service) GetDataBySn(orderSn string) (push *model.Push, err error) {

	push, err = s.dao.GetDataBySn(orderSn)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, orderSn, err)
	}

	return
}

//重启初始化--解锁
func (s *Service) InitPushLock() error {

	return s.dao.InitPushLock(LOCKTIMEOUT)
}
