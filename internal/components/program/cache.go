package program

import (
	"context"
	"errors"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/consts"
)

var (
	m     = sync.Map{} // service_name(routing_key):Programs
	empty = &Programs{}
)

func GetOrCreate(serviceName string) (*Programs, error) {
	if p, ok := m.Load(serviceName); ok {
		return p.(*Programs), nil
	}

	return buildCacheByService(serviceName)
}

func buildCacheByService(serviceName string) (*Programs, error) {
	sto := registry.Storages.GetStorage(consts.StorageNameProgram)
	kvs, err := sto.Get(context.Background(), serviceName)
	switch {
	case err == nil:
	case errors.Is(err, registry.ErrStorageNotFound):
		err = nil
	default:
		return nil, err
	}

	if len(kvs) == 0 {
		m.Store(serviceName, empty)
		return empty, nil
	}

	ps := &Programs{}
	for _, kv := range kvs {
		info := new(Info)
		if err := kv.Value.Scan(&info); err != nil {
			g.Log().Warningf(context.Background(), "scan program info failed: %v", err)
			continue
		}

		if err = ps.Create(info); err != nil {
			return nil, err
		}
	}
	m.Store(serviceName, ps)

	return ps, nil
}

func buildCache(ctx context.Context) {
	kvs, err := registry.Storages.GetStorage(consts.StorageNameProgram).Get(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "build program cache failed: %v", err)
		return
	}

	for _, kv := range kvs {
		info := new(Info)
		if err = kv.Value.Scan(&info); err != nil {
			g.Log().Warningf(ctx, "scan program info failed: %v", err)
			continue
		}
		var p *Programs
		if p, err = GetOrCreate(info.ServiceName); err != nil {
			g.Log().Errorf(ctx, "get or create programs failed: %v", err)
			continue
		}
		if err = p.Create(info); err != nil {
			g.Log().Errorf(ctx, "create program cache failed: %v", err)
			continue
		}
	}
}
