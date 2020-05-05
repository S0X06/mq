package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	TIMEOUT = (1 << 1)
	TTL     = (1 << 2)
)

type Register struct {
	grpcKey  string
	grpcAddr string
	httpKey  string
	httpUrl  string
	leaseid  clientv3.LeaseID
	client   *clientv3.Client
}

//创建对象
func NewRegister(addressLIst string) (*Register, error) {
	endpoints := strings.Split(addressLIst, ",")

	cli, err := clientv3.New(clientv3.Config{
		// 集群列表
		Endpoints:   endpoints,
		DialTimeout: TIMEOUT * time.Second,
	})

	if err != nil {
		return
	}

	return &Register{
		client: cli,
	}, err
}

//注册
func (this *Register) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), TTL)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), this.grpcKey, string(this.grpcAddr), clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, err
	}
	this.leaseid = resp.ID

	return this.client.KeepAlive(context.TODO(), resp.ID)

}

func (s *Register) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", s.Name)
	return err
}

//工作
func (this *Register) work() {

	ch, err := s.keepAlive()
	if err != nil {

	}

	for {
		select {
		case <-this.client.Ctx().Done():
			return errors.New("server closed")
		case ka, ok := <-ch:
			if !ok {
				continue
			}
			log.Printf("Recv reply from service: %s, ttl:%d", s.Name, ka.TTL)
		default:
			time.Sleep(5 * time.Second)

		}
	}

}
