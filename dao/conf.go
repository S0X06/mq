package dao

import (
	"errors"
	"mq/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//获取配置数据
func (d *Dao) GetConf() (*[]model.Conf, error) {
	conf := &[]model.Conf{}

	c, err := d.Collection(model.Conf{})
	if err != nil {
		return conf, err
	}

	err = c.Find(bson.M{"status": 1}).All(conf)
	return conf, err
}

//插入配置数据
func (d *Dao) InsertConf(confDate *model.Conf) (err error) {

	c, err := d.Collection(model.Conf{})
	if err != nil {
		return err
	}

	conf := &model.Conf{}

	err = c.Find(bson.M{"app_id": confDate.AppId}).One(conf)

	if err != mgo.ErrNotFound {
		return errors.New("服务已存在")
	}

	err = c.Insert(confDate)
	return
}

//更新发送状态
func (d *Dao) UpdateConf(confDate *model.Conf) error {

	c, err := d.Collection(model.Conf{})
	if err != nil {
		return err
	}

	conf := &model.Conf{}

	err = c.Find(bson.M{"app_id": confDate.AppId}).One(conf)
	if err != nil {
		return err
	}

	err = c.Update(bson.M{"app_id": confDate.AppId}, confDate)
	return err
}

//获取配置数据
func (d *Dao) GetRowConf(appId string) (*model.Conf, error) {
	conf := &model.Conf{}

	c, err := d.Collection(model.Conf{})
	if err != nil {
		return conf, err
	}
	err = c.Find(bson.M{"app_id": appId}).One(conf)
	return conf, err
}

//删除
func (d *Dao) RemoveConf(appId string) error {

	c, err := d.Collection(model.Conf{})
	if err != nil {
		return err
	}

	err = c.Remove(bson.M{"app_id": appId})
	return err
}
