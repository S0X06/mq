package utils

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	cHandler func(docs interface{}) (*mgo.Collection, error)

	hooker struct {
		c *mgo.Collection
	}

	//集合
	Log struct {
	}
)

//对外提供写函数
//获取mgo集合句柄
//文档结构体
func WriteLog(handler cHandler, data, message interface{}) {

	log := logrus.New()

	c, err := handler(&Log{})
	if err != nil {
		fmt.Println("日记集合获取失败")
		return
	}
	hooker, err := NewHooker(c)
	if err != nil {
		return
	}

	log.Hooks.Add(hooker)

	log.WithFields(logrus.Fields{
		"data": data,
	}).Error(message)

}

//初始化结构体
func NewHooker(c *mgo.Collection) (*hooker, error) {

	return &hooker{c: c}, nil
}

//实现接口函数
func (hk *hooker) Fire(entry *logrus.Entry) error {

	data := make(logrus.Fields)
	data["level"] = entry.Level.String()
	data["created_at"] = entry.Time
	data["message"] = entry.Message

	for k, v := range entry.Data {
		if errData, isError := v.(error); logrus.ErrorKey == k && v != nil && isError {
			data[k] = errData.Error()
		} else {
			data[k] = v
		}
	}

	mgoErr := hk.c.Insert(bson.M(data))

	if mgoErr != nil {
		return fmt.Errorf("mongodb: %v", mgoErr)
	}

	return nil
}

//实现接口函数
func (hk *hooker) Levels() []logrus.Level {
	return logrus.AllLevels
}
