package service

import (
	"fmt"

	"github.com/gregdel/pushover"
	"github.com/talbx/openwindow/pkg/model"
)

type Notifier interface {
	Notify(model.TuyaHumidity, NotificationType)
	buildMessage(tuya model.TuyaHumidity, t NotificationType) string
}

type NotificationService struct{}
type NotificationType int

const (
	FIRING   NotificationType = iota
	RESOLVED NotificationType = iota
)

func (n NotificationService) Notify(tuya model.TuyaHumidity, t NotificationType) {
	app := pushover.New(model.OWC.ApiToken)

	message := &pushover.Message{
		Message:  n.buildMessage(tuya, t),
		Title:    tuya.Device,
		Priority: pushover.PriorityHigh,
		Sound:    pushover.SoundCosmic,
	}
	recipient := pushover.NewRecipient(model.OWC.UserToken)
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		model.SugaredLogger.Panic(err)
	}

	model.SugaredLogger.Infof("Sent out pushover message with id %v", response.ID)
}

func (n NotificationService) buildMessage(tuya model.TuyaHumidity, t NotificationType) string {
	if t == FIRING {
		return fmt.Sprintf("%.2f humidity!", tuya.Humidity)
	}
	return fmt.Sprintf("Resolved! Humidity at %.2f okay again", tuya.Humidity)
}
