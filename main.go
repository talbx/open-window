package main

import (
	"os"
	"os/signal"
	"syscall"

	broker "github.com/talbx/openwindow/pkg/broker"
	model "github.com/talbx/openwindow/pkg/model"
)

var bc chan os.Signal

func main() {
	logger := model.CreateSugaredLogger()
	model.CreateOpenWindowConfig()
	bc = make(chan os.Signal, 1)
	signal.Notify(bc, os.Interrupt, syscall.SIGTERM)
	broker.Attach()
	sig := <-bc

	logger.Debugf("bye bye: %s", sig.String())
}