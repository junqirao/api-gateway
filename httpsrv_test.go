package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/grace"
	"api-gateway/internal/components/response"
	"api-gateway/internal/consts"
)

func TestHttpSrv0(t *testing.T) {
	startEchoServer("test0", 8997, 10)
}

func TestHttpSrv1(t *testing.T) {
	startEchoServer("test1", 8998, 10)
}

func TestHttpSrv2(t *testing.T) {
	startEchoServer("test2", 8999, 30)
}

func startEchoServer(name string, port, weight int) {
	id := fmt.Sprintf("server#%v", name)
	server := g.Server(id)
	server.SetPort(port)
	server.SetClientMaxBodySize(consts.DefaultMaxBodySize)
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/echo", func(r *ghttp.Request) {
			st := r.GetQuery("status", 200).Int()
			slp := r.GetQuery("sleep", 0).Int()
			if slp > 0 {
				time.Sleep(time.Millisecond * time.Duration(slp))
			}
			file := r.GetUploadFile("file")
			if file != nil {
				fmt.Printf("file: %s, size: %dbytes\n", file.Filename, file.Size)
				f, err := file.Open()
				if err != nil {
					response.WriteJSON(r, response.NewCode(st, id, 500).WithMessage(err.Error()))
					return
				}
				defer f.Close()
				local, err := gfile.OpenFile(fmt.Sprintf("./%s", file.Filename), os.O_CREATE|os.O_RDWR, gfile.DefaultPermCopy)
				if err != nil {
					response.WriteJSON(r, response.NewCode(st, id, 500).WithMessage(err.Error()))
					return
				}
				defer local.Close()
				n, err := io.Copy(local, f)
				if err != nil {
					response.WriteJSON(r, response.NewCode(st, id, 500).WithMessage(err.Error()))
					return
				}
				fmt.Printf("copy bytes: %d\n", n)
			}
			fmt.Printf("echo: %d, header:%+v\n", st, r.Header)
			response.WriteJSON(r, response.NewCode(st, id, st))
		})
		group.ALL("/registry", func(r *ghttp.Request) {
			services, err := registry.Registry.GetServices(r.Context())
			if err != nil {
				response.WriteJSON(r, response.NewCode(500, id, 500))
				return
			}
			m := make(map[string][]*registry.Instance)
			for k, v := range services {
				m[k] = v.Instances()
			}
			r.Response.WriteHeader(http.StatusOK)
			r.Response.WriteJson(m)
		})
	})
	cfg := registry.Config{}
	err := g.Cfg().MustGet(context.Background(), "registry").Struct(&cfg)
	if err != nil {
		panic(err)
		return
	}
	err = registry.Init(context.Background(), cfg, registry.NewInstance("test").WithMetaData(
		map[string]interface{}{
			"name":   name,
			"weight": weight,
		},
	).WithAddress("127.0.0.1", port))
	if err != nil {
		panic(err)
		return
	}
	grace.Register(context.Background(), "deregister_registry", func() {
		_ = registry.Registry.Deregister(context.Background())
	})
	server.Run()
	grace.ExecAndExit(context.Background())
}
