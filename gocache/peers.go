package gocache

import (
	"MiniArch/gocache/gcachepb"
	"context"
)

type PromoteGetterPicker interface {
	PickPromoteGetter(key string) (PromoteGetter, bool)
}
type PromoteGetter interface {
	GetFromPromote(ctx context.Context, req *gcachepb.GetRequest, resp *gcachepb.GetResponse) error
}
type NoPeers struct{}

func (n NoPeers) PickPromoteGetter(key string) (PromoteGetter, bool) {
	return nil, false
}
