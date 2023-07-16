package broker

import (
	"encoding/json"
	"errors"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/talbx/openwindow/pkg/model"
	"github.com/talbx/openwindow/pkg/service"
)

var received = false
var n = service.NotificationService{}
var change = service.ChangeService{N: n}

func Attach() {
	opts := createMqttOpts()
	opts.OnConnect = OnConnect
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		model.SugaredLogger.Error("could not etablish a stable connection to the broker using client::connect")
		model.SugaredLogger.Error(token.Error())
		os.Exit(1)
	} else {
		model.SugaredLogger.Infof("Connected to mosquitto instance on %v", model.OWC.MqttHost)
	}
}

func createMqttOpts() *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker(model.OWC.MqttHost)
	opts.SetClientID(model.OWC.MqttClientId)
	opts.SetDefaultPublishHandler(f)
	return opts
}

func OnConnect(c MQTT.Client) {
	for _, device := range model.OWC.Devices {
		if token := c.Subscribe(device.Topic, 0, f); token.Wait() && token.Error() != nil {
			model.SugaredLogger.Errorf("there was an error during topic subscription for %v", device.Topic)
			panic(token.Error())
		}
	}
	if received {
		println("receviced")
	}
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	var tuya model.TuyaHumidity
	json.Unmarshal(msg.Payload(), &tuya)
	room, err := translateName(msg.Topic())
	tuya.Device = room
	if err != nil {
		model.SugaredLogger.Errorf("the topic %v could not be translated into a room as defined in config.toml", err)
		return
	}
	change.HandleChange(tuya)
}

func translateName(topic string) (string, error) {
	for _, device := range model.OWC.Devices {
		if device.Topic == topic {
			return device.Room, nil
		}
	}
	return "", errors.New("There is no device configuration for this presented topic: " + topic)
}
