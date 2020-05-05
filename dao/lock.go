package dao

import (
	"mq/model"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//处理锁
func (d *Dao) ProcessLock(key string, value int, lockValue int) error {

	c, err := d.Collection(model.Lock{})
	if err != nil {
		return err
	}

	data := bson.M{
		"$set": bson.M{
			"value":      lockValue,
			"updated_at": time.Now().Unix(),
		},
	}

	err = c.Update(bson.M{"key": key, "value": value}, data)
	return err
}

//初始化key
func (d *Dao) UpdateLockKey(key string, value int, timeOut time.Duration) error {

	c, err := d.Collection(model.Lock{})
	if err != nil {
		return err
	}

	data := bson.M{
		"$set": bson.M{
			"value": 0,
		},
	}

	// err = c.Update(bson.M{"key": key}, data)
	c.Upsert(bson.M{"key": key}, data)
	return err
}

//初始化key
func (d *Dao) CheckLockKey(key string, value int, timeOut time.Duration) error {

	c, err := d.Collection(model.Lock{})
	if err != nil {
		return err
	}

	data := bson.M{
		"$set": bson.M{
			"value": value,
		},
	}

	beforeTime := time.Now().Add(time.Second * timeOut).Unix()
	// err = c.Update(bson.M{"key": key}, data)
	c.Update(bson.M{"key": key, "lock_at": bson.M{"$lte": beforeTime}}, data)
	return err
}

//查找
func (d *Dao) FindLock(key string) (*model.Lock, error) {

	lock := &model.Lock{}
	c, err := d.Collection(model.Lock{})
	if err != nil {
		return lock, err
	}

	err = c.Find(bson.M{"key": key}).One(lock)
	if err != nil {
		return lock, err
	}

	return lock, nil
}

//插入锁
func (d *Dao) InsertLock(lock *model.Lock) error {

	c, err := d.Collection(model.Lock{})
	if err != nil {
		return err
	}

	err = c.Insert(lock)
	return err
}
