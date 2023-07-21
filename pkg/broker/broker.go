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

type Broker struct {
	Change service.CanChange
	Exiter
}

func (b Broker) Attach() {
	opts := b.createMqttOpts()
	opts.OnConnect = b.OnConnect
	client := MQTT.NewClient(opts)
	token := client.Connect()
	b.Connect(token)
}

func (b Broker) Connect(token MQTT.Token) {
	if token.Wait() && token.Error() != nil {
		model.SugaredLogger.Error("could not etablish a stable connection to the broker using client::connect")
		b.Exiter.Exit(token.Error())
	} else {
		model.SugaredLogger.Infof("Connected to mosquitto instance on %v", model.OWC.MqttConfig.Host)
	}
}

func (e ExitHandler) Exit(tokenError error) {
	model.SugaredLogger.Error(tokenError)
	os.Exit(1)
}

func (b Broker) createMqttOpts() *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker(model.OWC.MqttConfig.Host)
	opts.SetClientID(model.OWC.MqttConfig.ClientId + "-" + time.Now().Format(time.DateOnly))
	model.SugaredLogger.Infof("Set the MQTT Client Id to %v", opts.ClientID)
	opts.SetDefaultPublishHandler(b.HandleMessage)
	return opts
}

func (b Broker) OnConnect(c MQTT.Client) {
	for _, device := range model.OWC.Devices {
		token := c.Subscribe(device.Topic, 0, b.HandleMessage)
		if token.Wait() && token.Error() != nil {
			model.SugaredLogger.Errorf("there was an error during topic subscription for %v", device.Topic)
			panic(token.Error())
		}
	}
}

func (b Broker) HandleMessage(_ MQTT.Client, msg MQTT.Message) {
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
	b.Change.HandleChange(tuya)
}

func Translate(topic string) (string, error) {
	for _, device := range model.OWC.Devices {
		if device.Topic == topic {
			return device.Room, nil
		}
	}
	return "", errors.New("There is no device configuration for this presented topic: " + topic)
}

type Exiter interface {
	Exit(error)
}

type ExitHandler struct{}
