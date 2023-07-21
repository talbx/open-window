package service

import (
	"time"

	"github.com/talbx/openwindow/pkg/model"
)

type NotifyBridge struct {
	RealNotifier Notifier
}

var NotifyCache = make(map[string]time.Time, 0)

func (bridge NotifyBridge) Notify(m model.TuyaHumidity, n NotificationType) {
	if n == RESOLVED {
		bridge.RealNotifier.Notify(m, n)
		return
	}

	lastSent, ok := NotifyCache[m.Device]
	if ok && n == FIRING {
		now := time.Now()
		diff := now.Sub(lastSent)
		if diff.Minutes() >= float64(model.GetGlobalConfig().Interval) {
			NotifyCache[m.Device] = now
			model.SugaredLogger.Infof("The NotifyBridge found out that the last notification for firing device %v is overdue by %v (>= 30 Minutes)", m.Device, diff.Minutes())
			bridge.RealNotifier.Notify(m, n)
			return
		}
		model.SugaredLogger.Infof("The NotifyBridge found out that the last notification for firing device %v is only %v minutes old. Therefore the notifier will not be triggered", m.Device, diff.Minutes())
	}

}
