package impl

import "fmt"
import "github.com/hashicorp/consul/api"
import "google.golang.org/grpc/naming"

type Service struct{}

func filterUnheathy(services []*api.ServiceEntry) []*api.ServiceEntry {
	var ret []*api.ServiceEntry
	for _, svc := range services {
		for _, check := range svc.Checks {
			if check.Status == "passing" && check.ServiceID == svc.Service.ID {
				ret = append(ret, svc)
			}
		}
	}
	return ret
}

func (*Service) Change(olds, news interface{}) []*naming.Update {
	fmt.Println("hello")
	var oldServices, newServices []*api.ServiceEntry
	if services, ok := olds.([]*api.ServiceEntry); ok {
		oldServices = services
	}

	if services, ok := news.([]*api.ServiceEntry); ok {
		newServices = services
	}

	var tmp = make(map[string]*naming.Update)
	oldServices = filterUnheathy(oldServices)
	newServices = filterUnheathy(newServices)
	// uid 防止不同机器上产生相同的容器ID
	for _, old := range oldServices {
		tmp[old.Service.ID] = &naming.Update{
			Op:   naming.Delete,
			Addr: fmt.Sprintf("%s:%d", old.Service.Address, old.Service.Port),
		}
	}

	for _, s := range newServices {
		if _, ok := tmp[s.Service.ID]; ok {
			delete(tmp, s.Service.ID)
		} else {
			// if heald
			tmp[s.Service.ID] = &naming.Update{
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

func (*Service) Param(target string) map[string]interface{} {
	return map[string]interface{}{"type": "service", "service": target}
}
