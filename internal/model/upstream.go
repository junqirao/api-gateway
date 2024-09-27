package model

import (
	"api-gateway/internal/components/upstream"
)

type GetServiceStateOutput struct {
	Detail map[string][]*upstream.UpsState
}
