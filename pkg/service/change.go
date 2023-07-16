package service

import (
	caching "github.com/talbx/openwindow/pkg/cache"
	"github.com/talbx/openwindow/pkg/model"
)

type ChangeService struct {
	N Notifier
}

var cache = caching.HumidityCache{Container: make(map[string]float32, 0)}

func (c ChangeService) HandleChange(h model.TuyaHumidity) {
	old := cache.Load(h.Device)
	cache.Store(h.Device, h.Humidity)

	if !c.IsOk(old) && c.IsOk(h.Humidity) {
		model.SugaredLogger.Infof("Humidity Resolved for %v; values in sweetspot again (%v). Will Notify!", h.Device, h.Humidity)
		c.N.Notify(h, RESOLVED)
		return
	} else if !c.IsOk(h.Humidity) {
		model.SugaredLogger.Infof("%s - have to notify since humidity is outside sweetspot (%.2f)", h.Device, h.Humidity)
		c.N.Notify(h, FIRING)
		return
	}
	model.SugaredLogger.Infof("%s - no notification sent, since humidity is in sweetspot (%.2f)", h.Device, h.Humidity)
}

func (c ChangeService) IsResolved(h model.TuyaHumidity) bool {
	hum := cache.Load(h.Device)
	if h.Humidity < hum {
		return c.IsOk(h.Humidity)
	}
	return c.IsOk(h.Humidity)
}

func (c ChangeService) IsOk(humidity float32) bool {
	return humidity == 0.0 || (humidity >= 40.0 && humidity < 60.0)
}
