package service

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"github.com/talbx/openwindow/pkg/model"
)

type FakeNotifier struct{
	Mock mock.Mock
}
var received model.TuyaHumidity
var receivedType NotificationType

func (f FakeNotifier) Notify(tuya model.TuyaHumidity, typus NotificationType){
	received = tuya
	receivedType = typus
}

func (f FakeNotifier) buildMessage(tuya model.TuyaHumidity, typus NotificationType) string{
	return ""
}

var F FakeNotifier = FakeNotifier{}

func Test_HandleChange(t *testing.T){
	model.CreateSugaredLogger()
	
	a1 := model.TuyaHumidity{
		Device: "DeviceA",
		Humidity: 62.3,
	}	
	a2 := model.TuyaHumidity{
		Device: "DeviceA",
		Humidity: 61.3,
	}

	b3 := model.TuyaHumidity{
		Device: "DeviceB",
		Humidity: 41.3,
	}

	a4 := model.TuyaHumidity{
		Device: "DeviceA",
		Humidity: 59.3,
	}

	service := ChangeService{N: &F}

	service.HandleChange(a1)
	assert.Equal(t, receivedType, FIRING)
	service.HandleChange(a2)
	assert.Equal(t, receivedType, FIRING)
	service.HandleChange(b3)
	F.Mock.AssertNotCalled(t,"Notify")
	service.HandleChange(a4)
	assert.Equal(t, receivedType, RESOLVED)

}