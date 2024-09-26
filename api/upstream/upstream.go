// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package upstream

import (
	"context"

	"api-gateway/api/upstream/v1"
)

type IUpstreamV1 interface {
	GetServiceNames(ctx context.Context, req *v1.GetServiceNamesReq) (res *v1.GetServiceNamesRes, err error)
	GetServiceState(ctx context.Context, req *v1.GetServiceStateReq) (res *v1.GetServiceStateRes, err error)
}
