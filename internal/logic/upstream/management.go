package upstream

import (
	"context"
	"sort"

	registry "github.com/junqirao/simple-registry"

	ups "api-gateway/api/inner/upstream"
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

	details := make(map[string][]*upstream.UpsState)
	for _, instance := range svc.Instances() {
		cc, err := utils.ClientConnFromInstance(instance)
		if err != nil {
			return nil, err
		}
		states, err := ups.NewManagementClient(cc).GetServiceStates(ctx, &ups.GetServiceStatesReq{ServiceName: serviceName})
		if err != nil {
			return nil, err
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
	}

	return &model.GetServiceStateOutput{
		Detail: details,
	}, nil
}
