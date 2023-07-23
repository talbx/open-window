package service

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/talbx/openwindow/pkg/model"
)

type FakeNotifier struct {
	Mock mock.Mock
}

var received model.TuyaHumidity
var receivedType NotificationType

func (f *FakeNotifier) Notify(tuya model.TuyaHumidity, typus NotificationType) {
	received = tuya
	receivedType = typus
}

func Test_HandleChange(t *testing.T) {
	model.CreateSugaredLogger()

	a1 := model.TuyaHumidity{
		Device:   "DeviceA",
		Humidity: 62.3,
	}
	a2 := model.TuyaHumidity{
		Device:   "DeviceA",
		Humidity: 61.3,
	}

	b3 := model.TuyaHumidity{
		Device:   "DeviceB",
		Humidity: 41.3,
	}

	a4 := model.TuyaHumidity{
		Device:   "DeviceA",
		Humidity: 59.3,
	}
	f := new(FakeNotifier)
	service := ChangeService{N: f}

	service.HandleChange(a1)
	assert.Equal(t, receivedType, FIRING)
	service.HandleChange(a2)
	assert.Equal(t, receivedType, FIRING)
	service.HandleChange(b3)
	f.Mock.AssertNotCalled(t, "Notify")
	service.HandleChange(a4)
	assert.Equal(t, receivedType, RESOLVED)

}

func Test_IsResolved(t *testing.T) {
	model.CreateSugaredLogger()

	a4 := model.TuyaHumidity{
		Device:   "DeviceA",
		Humidity: 59.3,
	}

	a5 := model.TuyaHumidity{
		Device:   "DeviceB",
		Humidity: 69.3,
	}

	f := new(FakeNotifier)
	service := ChangeService{N: f}
	storedHumidity["DeviceA"] = 61.5
	storedHumidity["DeviceB"] = 61.5
	resolved := service.IsResolved(a4)
	resolved2 := service.IsResolved(a5)

	assert.True(t, resolved)
	assert.False(t, resolved2)
}
