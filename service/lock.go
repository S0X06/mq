package service

//初始化
func (s *Service) InitLock() error {

	//初始化mq分布式锁
	// sendKey := "send"
	err := s.dao.CheckLockKey(SENDLOCK, 0, LOCKTIMEOUT)
	if err != nil {
		return err
	}

	//回调分布式锁
	// ackKey := "notify"
	err = s.dao.CheckLockKey(NOTIFYLOCK, 0, LOCKTIMEOUT)
	if err != nil {
		return err
	}

	return nil

}
