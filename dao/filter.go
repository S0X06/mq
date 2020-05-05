package dao

import (
	"mq/model"

	"gopkg.in/mgo.v2/bson"
)

func (d *Dao) BatchFilter(pushs *[]model.Push, ack int) error {

	filters, ok := model.ParsePushToFilter(pushs, ack)

	if !ok {
		return nil
	}

	c, err := d.Collection(model.Filter{})
	if err != nil {
		return err
	}

	err = c.Insert(filters...)
	return err

}

//创建一个过滤
func (d *Dao) InsertFilter(filter *model.Filter) error {

	c, err := d.Collection(model.Filter{})
	if err != nil {
		return err
	}

	err = c.Insert(filter)
	return err

}

//删除过滤缓存
func (d *Dao) RemoveFilter(orderSn string, ack int) error {

	c, err := d.Collection(model.Filter{})
	if err != nil {
		return err
	}

	err = c.Remove(bson.M{"order_sn": orderSn, "ack": ack})
	return err
}

//获取过滤
func (d *Dao) GetFilter(ack int) (*[]model.Filter, error) {

	filters := &[]model.Filter{}

	c, err := d.Collection(&model.Filter{})
	if err != nil {
		return filters, err
	}

	err = c.Find(bson.M{"ack": ack}).All(filters)
	return filters, err

}
