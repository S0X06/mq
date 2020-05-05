package dao

import (
	"errors"
	"fmt"
	"mq/conf"
	"mq/utils"
	"reflect"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	DBCHandle     func(docs interface{}) (*mgo.Collection, error)
	Lock          sync.RWMutex
	tableMapCache = make(map[string]string)
)

type Dao struct {
	conf      *conf.Config
	mogClient *MogClient
}

type MogClient struct {
	TmpStatus bool
	Status    chan bool
	Session   *mgo.Session
	Database  *mgo.Database
	// Lock      sync.RWMutex
}

func New(c *conf.Config) (dao *Dao) {

	dao = &Dao{
		conf: c,
		mogClient: &MogClient{
			TmpStatus: false,
			Status:    make(chan bool, 1),
		},
	}

	DBCHandle = dao.Collection
	dao.mogClient.Status <- false

	mgo.SetDebug(true)

	go dao.connect()

	return
}

func (d *Dao) Close() {
	defer d.mogClient.Session.Close()
}

func (d *Dao) connect() *mgo.Session {

	//配置
	conf := d.conf.MongoDB

	for {

		select {

		case status := <-d.mogClient.Status:

			//连接失败或未连接
			if !status {

				dial := &mgo.DialInfo{
					Addrs:  []string{conf.Addr},
					Direct: false,
					// Database:  conf.DdataBase,
					Timeout:   time.Second * 1,
					PoolLimit: 4096, // Session.SetPoolLimit
					// Username:  conf.UserName,
					// Password:  conf.PassWord,
				}
				//创建一个维护套接字池的session

				session, err := mgo.DialWithInfo(dial)
				if err != nil {
					fmt.Printf("MongoDB 连接失败,失败原因： %s。2 秒后 尝试重新连接....\n", err.Error())
					continue
				}

				session.SetMode(mgo.Monotonic, true)
				d.mogClient.Session = session

				//使用指定数据库
				database := session.DB(conf.DdataBase)

				d.mogClient.Database = database

				d.mogClient.TmpStatus = true
				d.mogClient.Status <- true

			} else {

				//ping
				if err := d.mogClient.Session.Ping(); err != nil {
					d.mogClient.TmpStatus = false
					continue
				}

				// //初始化
				// err := d.initPushLock()
				// if err != nil {
				// 	fmt.Printf("MongoDB push初始化失败 ：%v\n", err)
				// 	d.mogClient.TmpStatus = false
				// 	continue
				// }

				// //初始化
				// err = d.initLock()
				// if err != nil {
				// 	fmt.Printf("MongoDB lock初始化失败 ：%v\n", err)
				// 	d.mogClient.TmpStatus = false
				// 	continue
				// }

				fmt.Printf("MongoDB 连接成功 ...\n")

			}

		default:

			time.Sleep(2 * time.Second)
			if !d.mogClient.TmpStatus || d.mogClient.Session == nil || d.mogClient.Session.Ping() != nil {
				// fmt.Printf("mongoDB 重新连接\n")
				d.mogClient.TmpStatus = false
				d.mogClient.Status <- false
			}

		}

	}

}

//集合 == 表
func (d *Dao) Collection(docs interface{}) (*mgo.Collection, error) {

	if !d.mogClient.TmpStatus {
		return &mgo.Collection{}, errors.New("数据库未连接 ...")
	}

	var table string

	docsType := reflect.TypeOf(docs)
	docsType = utils.ParseReflect(docsType)
	name := docsType.Name()
	// fmt.Println(name)
	if _, ok := tableMapCache[name]; !ok {
		Lock.Lock()
		defer Lock.Unlock()

		table = strings.ToLower(name)
		tableMapCache[name] = table
	}

	table = tableMapCache[name]

	// fmt.Println(table)

	return d.mogClient.Database.C(table), nil
}
