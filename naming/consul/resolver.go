package consul

import "google.golang.org/grpc/naming"

type gRPCResolver struct {
	addr  string
	ttype int
}

func NewResolver(addr string, ttype int) naming.Resolver {
	return &gRPCResolver{
		addr:  addr,
		ttype: ttype,
	}
}

func (gr *gRPCResolver) Resolve(target string) (naming.Watcher, error) {
	w := &gRPCWatcher{
		addr:    gr.addr,
		target:  target,
		err:     make(chan error, 1),
		updates: make(chan []*naming.Update, 1),
		ttype:   gr.ttype,
	}
	return w, nil
}
