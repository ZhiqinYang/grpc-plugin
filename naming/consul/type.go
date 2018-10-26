package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc/naming"
)

type WatchType interface {
	Change(olds, news interface{}) []*naming.Update
	Param(string) map[string]interface{}
}

const (
	ServiceType = iota
)

var factory = map[int]WatchType{
	ServiceType: new(service),
}

// WathType implement
// service watch type
type service int

func (*service) Change(olds, news interface{}) []*naming.Update {
	var oldServices, newServices []*api.ServiceEntry
	if services, ok := olds.([]*api.ServiceEntry); ok {
		oldServices = services
	}

	if services, ok := news.([]*api.ServiceEntry); ok {
		newServices = services
	}

	var tmp = make(map[string]*naming.Update)
	// uid 防止不同机器上产生相同的容器ID
	for _, old := range oldServices {
		uid := fmt.Sprintf("%s-%s", old.Service.ID, old.Node.ID)
		tmp[uid] = &naming.Update{
			Op:   naming.Delete,
			Addr: fmt.Sprintf("%s:%d", old.Service.Address, old.Service.Port),
		}
	}

	for _, s := range newServices {
		uid := fmt.Sprintf("%s-%s", s.Service.ID, s.Node.ID)
		if _, ok := tmp[uid]; ok {
			delete(tmp, uid)
		} else {
			tmp[uid] = &naming.Update{
				Op:   naming.Add,
				Addr: fmt.Sprintf("%s:%d", s.Service.Address, s.Service.Port),
			}
		}
	}
	var ret = make([]*naming.Update, 0)
	for _, v := range tmp {
		ret = append(ret, v)
	}
	return ret
}

func (*service) Param(target string) map[string]interface{} {
	return map[string]interface{}{"type": "service", "service": target}
}
