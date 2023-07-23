package service

import (
	"fmt"
	"time"

	"github.com/talbx/openwindow/pkg/model"
)

type NotifyBridge struct {
	RealNotifier Notifier
}

var NotifyCache = make(map[string]time.Time, 0)

func (bridge NotifyBridge) Notify(m model.TuyaHumidity, n NotificationType) {
	model.SugaredLogger.Debugf("evaluating if %v with notification type '%v' needs to notify", m.Device, n)
	if n == RESOLVED {
		model.SugaredLogger.Debugf("%v is resolved, therefore will notify", m.Device)
		bridge.RealNotifier.Notify(m, n)
		delete(NotifyCache, m.Device)
		return
	}

	bridge.handleFiring(m)
}

func (bridge NotifyBridge) handleFiring(m model.TuyaHumidity) {

	lastSent, ok := NotifyCache[m.Device]
	fmt.Println(lastSent.String(), ok)
	if ok {
		bridge.handleFiredBefore(m, lastSent)
		return
	}
	bridge.handleNeverFiredBefore(m)
}

func (bridge NotifyBridge) handleFiredBefore(m model.TuyaHumidity, lastSent time.Time) {
	model.SugaredLogger.Debugf("%v was seen before with: %v", m.Device, lastSent.String())
	now := time.Now()
	diff := now.Sub(lastSent)
	if diff.Minutes() >= float64(model.GetGlobalConfig().Interval) {
		NotifyCache[m.Device] = now
		model.SugaredLogger.Infof("last notification for firing device %v is overdue by %v (>= 30 Minutes)", m.Device, diff.Minutes())
		bridge.RealNotifier.Notify(m, FIRING)
		return
	}
	model.SugaredLogger.Debugf("last notification for firing device %v is only %v minutes old. Therefore the notifier will not be triggered", m.Device, diff.Minutes())
}

func (bridge NotifyBridge) handleNeverFiredBefore(m model.TuyaHumidity) {
	NotifyCache[m.Device] = time.Now()
	model.SugaredLogger.Infof("there was no sent information before for firing device %v. Therefore will notify", m.Device)
	bridge.RealNotifier.Notify(m, FIRING)
}
