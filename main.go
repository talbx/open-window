package main

import (
	"github.com/gregdel/pushover"
	"github.com/talbx/openwindow/pkg/service"
	"os"
	"os/signal"
	"syscall"

	broker "github.com/talbx/openwindow/pkg/broker"
	model "github.com/talbx/openwindow/pkg/model"
)

var bc chan os.Signal

func main() {
	logger := model.CreateSugaredLogger()
	logger.Info("open-window started")
	model.CreateOpenWindowConfig()
	bc = make(chan os.Signal, 1)
	signal.Notify(bc, os.Interrupt, syscall.SIGTERM)
	var n = service.NotifyBridge{RealNotifier: service.NotificationService{App: pushover.New(model.OWC.PushoverConfig.ApiToken)}}
	var change = service.ChangeService{N: n}
	b := broker.Broker{Change: change, Exiter: broker.ExitHandler{}}
	b.Attach()
	sig := <-bc
	logger.Debugf("open-window shutting down with signal %v", sig.String())
}
