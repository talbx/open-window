package service

import "github.com/talbx/openwindow/pkg/model"

type ChangeService struct{}

var storedHumidity *model.TuyaHumidity = nil

func (c ChangeService) IsResolved(h model.TuyaHumidity) bool {
	if h.Humidity < c.LoadStoredHumidity().Humidity {
		return c.IsOk(h.Humidity)
	}
	return false
}

func (c ChangeService) StoreHumidity(h model.TuyaHumidity) {
	model.SugaredLogger.Debugf("Stored humidity %v in cache", h)
	storedHumidity = &h
}

func (c ChangeService) LoadStoredHumidity() model.TuyaHumidity {
	model.SugaredLogger.Debugf("Loaded stored humidity %v from cache", storedHumidity)
	return *storedHumidity
}

func (c ChangeService) IsOk(humidity float32) bool {
	return humidity >= 40.0 && humidity <= 60.0
}