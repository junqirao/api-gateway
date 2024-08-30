package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/response"
)

func TestHttpSrv0(t *testing.T) {
	startEchoServer("test0", 8997)
}

func TestHttpSrv1(t *testing.T) {
	startEchoServer("test1", 8998)
}

func TestHttpSrv2(t *testing.T) {
	startEchoServer("test2", 8999)
}

func startEchoServer(name string, port int) {
	id := fmt.Sprintf("server#%v", name)
	server := g.Server(id)
	server.SetPort(port)
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/echo", func(r *ghttp.Request) {
			st := r.GetQuery("status", 200).Int()
			fmt.Printf("echo: %d\n", st)
			response.WriteJSON(r, response.NewCode(st, id, st))
		})
	})
	cfg := registry.Config{}
	err := g.Cfg().MustGet(context.Background(), "registry").Struct(&cfg)
	if err != nil {
		panic(err)
		return
	}
	err = registry.Init(context.Background(), cfg, registry.NewInstance("test").WithAddress("127.0.0.1", port))
	if err != nil {
		panic(err)
		return
	}
	server.Run()
}
