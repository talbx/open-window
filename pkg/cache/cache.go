package cache

import (
	"time"

	"github.com/talbx/openwindow/pkg/model"
)

type Cache interface {
	Store(key string, obj float32)
	Load(key string) float32
}

type HumidityCache struct{
	Container map[string]HumidityLog
}

type HumidityLog struct {
	Room string
	Hum float32
}

func (c HumidityCache) Store(device string, humidity float32) {
	model.SugaredLogger.Debugf("Stored humidity %v in cache for device %v", humidity, device)
	now := time.Now().String()
	c.Container[now] = HumidityLog{
		device,
		humidity,
	}
}

func (c HumidityCache) Load(device string) float32 {
	h := c.Container[device]
	model.SugaredLogger.Debugf("Loaded %v's stored humidity %v from cache", device, h)
	return h
}
