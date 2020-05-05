package dao

import (
	"mq/model"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//获取条件获取所有数据
func (d *Dao) GetPush(lock, ack, limit int) (*[]model.Push, error) {
	pushs := &[]model.Push{}

	c, err := d.Collection(model.Push{})
	if err != nil {
		return pushs, err
	}

	err = c.Find(bson.M{"ack": ack, "lock": 0}).Sort("sort").Limit(limit).All(pushs)
	// err = c.Find(bson.M{"ack": ack}).Limit(limit).All(pushs)
	return pushs, err
}

// select * from table where unix_timestamp(Time) > unix_timestamp('2010-04-11 22:55:00') and unix_timestamp(Time)< unix_timestamp('2010-04-11 23:00:00');
func (d *Dao) GetTimeOutAck(ack, lock, limit int, timeOut time.Duration) (*[]model.Push, error) {

	pushs := &[]model.Push{}

	c, err := d.Collection(model.Push{})
	if err != nil {
		return pushs, err
	}

	beforeTime := time.Now().Add(time.Second * timeOut).Unix()
	err = c.Find(bson.M{"ack": ack, "lock": lock, "created_at": bson.M{"$lte": beforeTime}}).Sort("sort").Limit(limit).All(pushs)

	return pushs, err
}

//通过获取消息
func (d *Dao) GetDataBySn(orderSn string) (*model.Push, error) {
	push := &model.Push{}

	c, err := d.Collection(model.Push{})
	if err != nil {
		return push, err
	}

	err = c.Find(bson.M{"order_sn": orderSn}).One(push)
	return push, err
}

//创建一个消息
func (d *Dao) CreatePush(push *model.Push) error {

	c, err := d.Collection(model.Push{})
	if err != nil {
		return err
	}

	err = c.Insert(push)
	return err
}

//更新发送状态
func (d *Dao) UpdateAck(orderSn string, ack int32, setAck int32) error {

	c, err := d.Collection(model.Push{})
	if err != nil {
		return err
	}

	push := &model.Push{
		// Ack: setAck,
	}

	err = c.Find(bson.M{"order_sn": orderSn}).One(push)
	if err != nil {
		return err
	}

	push.Ack = setAck

	err = c.Update(bson.M{"order_sn": orderSn, "ack": ack}, push)
	return err
}

//通用用更新
func (d *Dao) Set(where, data bson.M) error {

	c, err := d.Collection(model.Push{})
	if err != nil {
		return err
	}

	err = c.Update(where, data)
	return err
}

//批量更新
func (d *Dao) BatchUpdateAck(orderSns []string, push *model.Push) error {

	c, err := d.Collection(model.Push{})
	if err != nil {
		return err
	}

	err = c.Update(bson.M{"order_sn": bson.M{"$in": orderSns}}, &push)
	return err
}

//删除
func (d *Dao) RemovePush(orderSn string) error {

	c, err := d.Collection(model.Push{})
	if err != nil {
		return err
	}

	err = c.Remove(bson.M{"order_sn": orderSn})
	return err
}

//批量上锁
func (d *Dao) BatchLock(orderSns []string, lock int) error {

	c, err := d.Collection(model.Push{})
	if err != nil {
		return err
	}

	data := bson.M{
		"$set": bson.M{
			"lock":    lock,
			"lock_at": time.Now().Unix(),
		},
	}

	_, err = c.UpdateAll(bson.M{"order_sn": bson.M{"$in": orderSns}}, &data)
	return err
}

//重启初始化--解锁
func (d *Dao) InitPushLock(timeOut time.Duration) error {

	c, err := d.Collection(model.Push{})
	if err != nil {
		return err
	}

	data := bson.M{
		"$set": bson.M{
			"lock": 0,
		},
	}

	beforeTime := time.Now().Add(time.Second * timeOut).Unix()

	_, err = c.UpdateAll(bson.M{"lock": 1, "lock_at": bson.M{"$lte": beforeTime}}, &data)
	// fmt.Println(info)
	return err
}
