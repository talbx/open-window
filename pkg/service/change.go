package service

import (
	"github.com/talbx/openwindow/pkg/model"
)

type ChangeService struct {
	N Notifier
}

var storedHumidity = make(map[string]float32)

func (c ChangeService) HandleChange(h model.TuyaHumidity) {
	old := c.LoadStoredHumidity(h.Device)
	c.StoreHumidity(h.Device, h.Humidity)

	if !c.IsOk(old) && c.IsOk(h.Humidity) {
		model.SugaredLogger.Infof("Humidity Resolved for %v; values in sweetspot again (%v). Will Notify!", h.Device, h.Humidity)
		c.N.Notify(h, RESOLVED)
		return
	} else if !c.IsOk(h.Humidity) {
		model.SugaredLogger.Infof("%s - Humidity is outside sweetspot (%.2f)", h.Device, h.Humidity)
		c.N.Notify(h, FIRING)
		return
	}
	model.SugaredLogger.Infof("%s - no notification needed, since humidity is in sweetspot (%.2f)", h.Device, h.Humidity)
}

func (c ChangeService) IsResolved(h model.TuyaHumidity) bool {
	hum := c.LoadStoredHumidity(h.Device)
	if h.Humidity < hum {
		return c.IsOk(h.Humidity)
	}
	return c.IsOk(h.Humidity)
}

func (c ChangeService) IsOk(humidity float32) bool {
	return humidity == 0.0 || (humidity >= 40.0 && humidity < 60.0)
}

func (c ChangeService) StoreHumidity(device string, humidity float32) {
	model.SugaredLogger.Debugf("Stored humidity %v in cache for device %v", humidity, device)
	storedHumidity[device] = humidity
}

func (c ChangeService) LoadStoredHumidity(device string) float32 {
	h := storedHumidity[device]
	model.SugaredLogger.Debugf("Loaded %v's stored humidity %v from cache", device, h)
	return h
}
