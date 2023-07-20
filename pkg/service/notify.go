package service

import (
	"fmt"

	"github.com/gregdel/pushover"
	"github.com/talbx/openwindow/pkg/model"
)

type Notifier interface {
	Notify(model.TuyaHumidity, NotificationType)
}

var _ Notifier = NotificationService{}

type Meta struct {
	Priority int
	Sound    string
}
type NotificationService struct{}
type NotificationType int

const (
	FIRING   NotificationType = iota
	RESOLVED NotificationType = iota
)

func (n NotificationService) Notify(tuya model.TuyaHumidity, t NotificationType) {
	app := pushover.New(model.OWC.ApiToken)

	meta := buildMeta(t)
	message := &pushover.Message{
		Message:  n.buildMessage(tuya, t),
		Title:    tuya.Device,
		Priority: meta.Priority,
		Sound:    meta.Sound,
	}
	recipient := pushover.NewRecipient(model.OWC.UserToken)
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		model.SugaredLogger.Error(err)
		model.SugaredLogger.Errorf("No pushover message sent out due to an error communicating with the pusover api!")
		return
	}

	model.SugaredLogger.Infof("Sucessfully sent out pushover message with id %v", response.ID)
}

func (n NotificationService) buildMessage(tuya model.TuyaHumidity, t NotificationType) string {
	if t == FIRING {
		return fmt.Sprintf("%.2f humidity!", tuya.Humidity)
	}
	return fmt.Sprintf("Resolved! Humidity at %.2f okay again", tuya.Humidity)
}

func buildMeta(t NotificationType) Meta {
	if t == FIRING {
		return Meta{
			Priority: pushover.PriorityHigh,
			Sound:    pushover.SoundCosmic,
		}
	}
	return Meta{
		Priority: pushover.PriorityLow,
		Sound:    pushover.SoundMagic,
	}
}
