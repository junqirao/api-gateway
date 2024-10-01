package ipgeo

import (
	"context"
	"testing"
)

func TestIPGEO(t *testing.T) {
	Init(context.Background())

	t.Log(country("8.8.8.8"))
	t.Log(country("8.8.4.4"))
	t.Log(city("8.8.8.8"))
	t.Log(city("223.5.5.5"))
	t.Log(cityEN("8.8.8.8"))
	t.Log(cityEN("223.5.5.5"))
}
