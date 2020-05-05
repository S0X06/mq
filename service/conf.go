package service

import (
	"mq/dao"
	"mq/model"
	. "mq/queue"
	"mq/utils"
	"time"
)

//获取服务配置
func (s *Service) InitConf() error {
	//插入测试数据

	conf := &model.Conf{
		AppId:        "123",
		AppSecret:    "test",
		ServceName:   "kkk",
		QueueName:    "demo",
		RoutingKey:   "demo",
		ExchangeName: "demo",
		ExchangeType: "topic",
		Status:       1,
		UpdatedAt:    time.Now().Unix(),
		CreatedAt:    time.Now().Unix(),
	}

	s.dao.InsertConf(conf)

	//获取配置
	return nil
}

//获取配置
func (s *Service) GetConf(appId string) (model.Conf, error) {

	// SrvHandle.InitConf()

	if _, ok := ServiceConf[appId]; !ok {

		// fmt.Println("获取服务")
		conf, err := s.dao.GetRowConf(appId)
		if err != nil {
			utils.WriteLog(dao.DBCHandle, appId, err)
			return model.Conf{}, ErrNotFound
		}

		//判断是否发布
		if conf.Status != 1 {
			utils.WriteLog(dao.DBCHandle, appId, "服务未发布")
			return model.Conf{}, ErrNotFound
		}

		//加入配置
		Lock.Lock()
		defer Lock.Unlock()

		ServiceConf[appId] = *conf
		//加入通道
		Sdl.Conf <- *conf
	}
	//返回
	// fmt.Println(ServiceConf)
	return ServiceConf[appId], nil

}

//获取配置
func (s *Service) GetAllConf() (conf *[]model.Conf, err error) {
	conf, err = s.dao.GetConf()
	return
}

//添加配置
func (s *Service) CreateConf(conf *model.Conf) error {

	conf.Status = 0
	conf.CreatedAt = time.Now().Unix()
	conf.UpdatedAt = time.Now().Unix()

	err := s.dao.InsertConf(conf)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, conf, err)
	}
	return err
}

//更新
func (s *Service) UpdateConf(conf *model.Conf) error {

	err := s.dao.UpdateConf(conf)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, conf, err)
	}

	//从缓存中移除
	delete(ServiceConf, conf.AppId)

	return err
}

//删除
func (s *Service) RemoveConf(appId string) error {

	err := s.dao.RemoveConf(appId)
	if err != nil {
		utils.WriteLog(dao.DBCHandle, appId, err)
		return err
	}

	//从缓存中移除
	delete(ServiceConf, appId)

	return nil
}
