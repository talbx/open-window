package service

import (
	"github.com/stretchr/testify/mock"
	model2 "github.com/talbx/openwindow/pkg/model"
	"testing"
	"time"
)

type notifyMock struct {
	mock.Mock
	Notifier
}

func Test_NoNotificationForRecentEntry(t *testing.T) {

	model2.CreateSugaredLogger()
	mock := new(notifyMock)
	mock.On("Notify")
	bridge := NotifyBridge{RealNotifier: mock}

	NotifyCache["A"] = time.Now()

	model := model2.TuyaHumidity{Humidity: 60, Device: "A"}
	bridge.Notify(model, FIRING)

	mock.AssertNotCalled(t, "Notify")
}

func Test_NotificationForUnRecentEntry(t *testing.T) {

	model2.CreateSugaredLogger()
	mock := new(notifyMock)
	mock.On("Notify")
	bridge := NotifyBridge{RealNotifier: mock}

	NotifyCache["A"] = time.Now().Add(-31 * time.Minute)

	model := model2.TuyaHumidity{Humidity: 60, Device: "A"}
	bridge.Notify(model, FIRING)

	mock.AssertNumberOfCalls(t, "Notify", 1)
}

func Test_NotificationForUnRecentEntry2(t *testing.T) {

	model2.CreateSugaredLogger()
	mock := new(notifyMock)
	mock.On("Notify")
	bridge := NotifyBridge{RealNotifier: mock}

	NotifyCache["A"] = time.Now().Add(-61 * time.Minute)
	model := model2.TuyaHumidity{Humidity: 60, Device: "A"}
	bridge.Notify(model, FIRING)
	bridge.Notify(model, FIRING)
	bridge.Notify(model, FIRING)
	bridge.Notify(model, FIRING)

	NotifyCache["A"] = time.Now().Add(-31 * time.Minute)
	bridge.Notify(model, FIRING)

	mock.AssertNumberOfCalls(t, "Notify", 2)
}

func Test_Notify_Resolved(t *testing.T) {

	model2.CreateSugaredLogger()
	m := new(notifyMock)
	m.On("Notify")
	bridge := NotifyBridge{RealNotifier: m}

	model := model2.TuyaHumidity{Humidity: 30, Device: "A"}
	bridge.Notify(model, RESOLVED)
	m.AssertNumberOfCalls(t, "Notify", 1)
}

func (m *notifyMock) Notify(model2.TuyaHumidity, NotificationType) {
	m.Called()
}

func prepareMap() map[string]time.Time {
	m := make(map[string]time.Time, 0)
	m["A"] = time.Now()
	return nil
}
