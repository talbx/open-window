package main

import (
	"os"
	"os/signal"
	"syscall"

	broker "github.com/talbx/openwindow/pkg/broker"
	model "github.com/talbx/openwindow/pkg/model"
)

var c chan os.Signal

func main() {
	logger := model.CreateSugaredLogger()
	model.CreateOpenWindowConfig()
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	broker.Attach()
	ch := <-c

	logger.Debugf("bye bye: %s", ch.String())
}