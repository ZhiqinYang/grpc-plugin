package consul

import "google.golang.org/grpc/naming"
import "grpc-plugin/naming/consul/impl"

const (
	_ = iota
	ServiceBuilder
)

func init() {
	Registry(ServiceBuilder, &impl.Service{})
}

type WatchBuilder interface {
	Change(olds, news interface{}) []*naming.Update
	Param(string) map[string]interface{}
}

var factory = map[int]WatchBuilder{}

func Registry(ttype int, wb WatchBuilder) {
	factory[ttype] = wb
}
