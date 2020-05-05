package service

import (
	"sync"
	// "fmt"
	"errors"
	"mq/conf"
	"mq/dao"
	"mq/model"
	. "mq/queue"
)

const (
	//回调协议
	GRPC = (1 << 0)
	HTTP = (1 << 1)

	//ack 标志
	REMOVEACK = (0 << 1) //0
	PREACK    = (1 << 0) //1
	ACK       = (1 << 1) //2

	//push 锁
	LOCK   = (1 << 0)
	UNLOCK = (0 << 1)

	//模拟分布式抢锁
	SENDLOCK   = "send"
	NOTIFYLOCK = "notify"

	//排序
	ZEROSORT = (0 << 1)
	INCSORT  = (1 << 0)

	//超时回调
	TIMEOUT = (^(1 << 0)) //两秒

	LOCKTIMEOUT = (^(1 << 2)) //

	//查询数目
	LIMIT = (1 << 10)

	//回调函数标志
	OP_REMOVE = "OpPushRemove"
	OP_UPDATE = "OpPushUpdate"
	OP_SORT   = "OpPushSort"
	OP_UnLock = "OpPushUnLock"
)

var (
	SrvHandle *Service
	Lock      sync.RWMutex

	//全局数据
	SendPushs   *[]model.Push
	RequstPushs *[]model.Push
	ServiceConf = make(map[string]model.Conf)
	MapFn       = make(map[string]func(orderSn string))
	ProcessLock func(string, int, int) error

	ErrNotFound = errors.New("not found")
	ErrFound    = errors.New("found")
)

type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

func Init(c *conf.Config) (err error) {

	SrvHandle = New(c)

	//锁
	ProcessLock = SrvHandle.dao.ProcessLock

	//数据回调
	MapFn = map[string]func(orderSn string){
		OP_REMOVE: SrvHandle.OpPushRemove,
		OP_UPDATE: SrvHandle.OpPushUpdate,
		OP_SORT:   SrvHandle.OpPushSort,
		OP_UnLock: SrvHandle.OpPushUnLock,
	}

	//拷贝回调
	Sdl.TmpMapCallBack <- MapFn

	return
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}

	return
}

func (s *Service) Close() {
	s.dao.Close()
}
