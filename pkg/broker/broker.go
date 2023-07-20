package broker

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/talbx/openwindow/pkg/model"
	"github.com/talbx/openwindow/pkg/service"
)

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
		model.SugaredLogger.Infof("Connected to mosquitto instance on %v", model.OWC.MqttConfig.Host)
	}
}

func createMqttOpts() *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker(model.OWC.MqttConfig.Host)
	opts.SetClientID(model.OWC.MqttConfig.ClientId + "-" + time.Now().Format(time.DateOnly))
	model.SugaredLogger.Infof("Set the MQTT Client Id to %v", opts.ClientID)
	opts.SetDefaultPublishHandler(handleMessage)
	return opts
}

func OnConnect(c MQTT.Client) {
	for _, device := range model.OWC.Devices {
		if token := c.Subscribe(device.Topic, 0, handleMessage); token.Wait() && token.Error() != nil {
			model.SugaredLogger.Errorf("there was an error during topic subscription for %v", device.Topic)
			panic(token.Error())
		}
	}
}

func handleMessage(_ MQTT.Client, msg MQTT.Message) {
	var tuya model.TuyaHumidity
	err := json.Unmarshal(msg.Payload(), &tuya)

	if err != nil {
		model.SugaredLogger.Errorf("the message payload could not be unmarshaled: %v", err)
	}

	room, err := Translate(msg.Topic())
	tuya.Device = room
	if err != nil {
		model.SugaredLogger.Errorf("the topic %v could not be translated into a room as defined in config.toml", err)
		return
	}
	change.HandleChange(tuya)
}

func Translate(topic string) (string, error) {
	for _, device := range model.OWC.Devices {
		if device.Topic == topic {
			return device.Room, nil
		}
	}
	return "", errors.New("There is no device configuration for this presented topic: " + topic)
}
