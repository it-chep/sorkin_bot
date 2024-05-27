package adapter

import (
	"sorkin_bot/pkg/client/inmemory_cache"
	"time"
)

type AppointmentServiceAdapter struct {
	cache   inmemory_cache.Cache[string, any]
	gateway Gateway
}

func NewAppointmentServiceAdapter(gateway Gateway) *AppointmentServiceAdapter {
	return &AppointmentServiceAdapter{
		cache:   *inmemory_cache.NewCache[string, any](time.Second * 10),
		gateway: gateway,
	}
}
