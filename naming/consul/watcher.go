package consul

// cuonsul nameing with service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/consul/watch"

	"google.golang.org/grpc/naming"
)

type gRPCWatcher struct {
	addr    string
	target  string
	ttype   int
	plan    *watch.Plan
	err     chan error
	updates chan []*naming.Update
	once    sync.Once
}

func (gw *gRPCWatcher) Next() ([]*naming.Update, error) {
	gw.once.Do(func() {
		gw.watch()
	})
	select {
	case err := <-gw.err:
		return nil, err
	case updates := <-gw.updates:
		return updates, nil
	}
}

var (
	UnknowWathTypeErr = errors.New("unknow watch type")
)

func (gw *gRPCWatcher) watch() {
	wt, ok := factory[gw.ttype]
	if !ok {
		gw.err <- UnknowWathTypeErr
		return
	}

	plan, err := watch.Parse(wt.Param(gw.target))
	if err != nil {
		gw.err <- err
		return
	}

	var olds interface{}
	// 维护删除更新逻辑
	plan.Handler = func(idx uint64, raw interface{}) {

		fmt.Println("idx = ", idx, raw)

		if raw == nil {
			return
		}
		gw.updates <- wt.Change(olds, raw)
		olds = raw
	}

	go func() {
		if err := plan.Run(gw.addr); err != nil {
			gw.err <- err
		}
	}()
	gw.plan = plan
}

func (gw *gRPCWatcher) Close() {
	if gw.plan != nil {
		gw.plan.Stop()
	}
}
