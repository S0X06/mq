package etcd

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
)

var cli *clientv3.Client

type Resolver struct {
	rawAddr string
	cc      resolver.ClientConn
}

func NewResolver(etcdAddr string) resolver.Builder {
	return &etcdResolver{rawAddr: etcdAddr}
}
