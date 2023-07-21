package service

import (
	"github.com/gregdel/pushover"
	"github.com/stretchr/testify/mock"
	model2 "github.com/talbx/openwindow/pkg/model"
	"testing"
)

type pushoverMock struct {
	mock.Mock
}

func (p *pushoverMock) SendMessage(message *pushover.Message, recipient *pushover.Recipient) (*pushover.Response, error) {
	p.Called(message, recipient)
	return &pushover.Response{
		Status:  200,
		ID:      "Success",
		Receipt: "It worked",
	}, nil
}
func TestNotify_Firing(t *testing.T) {
	// given
	model2.OWC.PushoverConfig.UserToken = "ABC"
	model2.CreateSugaredLogger()
	msg, rec := buildArgs(pushover.SoundCosmic, pushover.PriorityHigh, "Some Device", "60.00 humidity!", "ABC")
	model := model2.TuyaHumidity{Humidity: 60.0, Device: "Some Device"}
	var m = new(pushoverMock)
	n := NotificationService{App: m}

	// when
	m.On("SendMessage", msg, rec)
	n.Notify(model, FIRING)

	// then
	m.AssertCalled(t, "SendMessage", msg, rec)
}
func TestNotify_Resolved(t *testing.T) {

	// given
	model2.OWC.PushoverConfig.UserToken = "ABC"
	model2.CreateSugaredLogger()
	msg, rec := buildArgs(pushover.SoundMagic, pushover.PriorityLow, "Some Device", "Resolved! Humidity at 59.00 okay again", "ABC")
	model := model2.TuyaHumidity{Humidity: 59.0, Device: "Some Device"}
	var m = new(pushoverMock)
	n := NotificationService{App: m}

	// when
	m.On("SendMessage", msg, rec)
	n.Notify(model, RESOLVED)

	// then
	m.AssertCalled(t, "SendMessage", msg, rec)
}

func buildArgs(sound string, priority int, device string, message string, recipient string) (*pushover.Message, *pushover.Recipient) {
	return &pushover.Message{Message: message, Title: device, Sound: sound, Priority: priority},
		pushover.NewRecipient(recipient)
}
