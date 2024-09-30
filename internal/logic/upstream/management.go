package upstream

import (
	"context"
	"sort"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	registry "github.com/junqirao/simple-registry"
	"google.golang.org/grpc"

	ups "api-gateway/api/inner/upstream"
	"api-gateway/internal/components/authentication"
	r "api-gateway/internal/components/registry"
	"api-gateway/internal/components/upstream"
	"api-gateway/internal/components/utils"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

type (
	sUpstreamManagement struct{}
)

func init() {
	service.RegisterUpstreamManagement(&sUpstreamManagement{})
}

func (s sUpstreamManagement) GetServiceNames(_ context.Context) []string {
	var ss []string
	upstream.Cache.Range(func(key, _ any) bool {
		ss = append(ss, key.(string))
		return true
	})
	sort.Strings(ss)
	return ss
}

func (s sUpstreamManagement) GetServiceState(ctx context.Context, serviceName string) (output *model.GetServiceStateOutput, err error) {
	svc, err := registry.Registry.GetService(ctx)
	if err != nil {
		return
	}
	output = new(model.GetServiceStateOutput)
	output.Detail = make(map[string][]any)
	getServiceStatus := func(ctx context.Context,
		cc *grpc.ClientConn,
		instance *registry.Instance) {

		ctx, cancelFunc := context.WithTimeout(ctx, time.Second*1)
		defer cancelFunc()
		states, err := ups.NewManagementClient(cc).GetServiceStates(ctx, &ups.GetServiceStatesReq{
			ServiceName:    serviceName,
			InstanceId:     r.CurrentInstance.Id,
			Authentication: authentication.L.Encode(r.CurrentInstance.Id),
		})
		if err != nil {
			output.Detail[instance.HostName] = append(output.Detail[instance.HostName], map[string]any{
				"error": err.Error(),
			})
			return
		}
		for _, state := range states.States {
			output.Detail[instance.HostName] = append(output.Detail[instance.HostName], &upstream.UpsState{
				HostName:     state.GetHostname(),
				InstanceId:   state.GetInstanceId(),
				Healthy:      state.GetHealthy(),
				Weight:       state.GetWeight(),
				WeightRatio:  float64(state.GetWeightRatio()),
				Load:         state.GetLoad(),
				BreakerState: state.GetBreakerState(),
			})
		}
		return
	}

	var cc *grpc.ClientConn
	for _, instance := range svc.Instances() {
		cc, err = utils.ClientConnFromInstance(instance, grpcx.Client.DefaultGrpcDialOptions()...)
		if err != nil {
			return
		}
		getServiceStatus(ctx, cc, instance)
	}

	return
}

func (s sUpstreamManagement) getServiceStatus(ctx context.Context,
	cc *grpc.ClientConn,
	serviceName string,
	details map[string][]*upstream.UpsState,
	instance registry.Instance) (err error) {
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*2)
	defer cancelFunc()
	states, err := ups.NewManagementClient(cc).GetServiceStates(ctx, &ups.GetServiceStatesReq{ServiceName: serviceName})
	if err != nil {
		return
	}
	for _, state := range states.States {
		details[instance.HostName] = append(details[instance.HostName], &upstream.UpsState{
			HostName:     state.GetHostname(),
			InstanceId:   state.GetInstanceId(),
			Healthy:      state.GetHealthy(),
			Weight:       state.GetWeight(),
			WeightRatio:  float64(state.GetWeightRatio()),
			Load:         state.GetLoad(),
			BreakerState: state.GetBreakerState(),
		})
	}
	return nil
}
