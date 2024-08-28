package config

import (
	"context"
)

func Init(ctx context.Context) {
	loadConfigs(ctx)
}
