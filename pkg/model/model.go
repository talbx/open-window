package model

import (
	"os"

	"github.com/BurntSushi/toml"
)

type TuyaHumidity struct {
	Device      string  `json:"device"`
	Battery     int     `json:"battery"`
	Humidity    float32 `json:"humidity"`
	Linkquality int     `json:"linkquality"`
	Temperature float32 `json:"temperature"`
	Voltage     int     `json:"voltage"`
}

type Device struct {
	Topic string
	Room string
}

type OpenWindowConfig struct {
	ApiToken     string
	UserToken    string
	MqttHost     string `toml:"host"`
	MqttClientId string `toml:"clientId"`
	Devices      []Device `toml:"devices"`
}

var OWC OpenWindowConfig

func CreateOpenWindowConfig(){
	file, err := os.ReadFile("config.toml")
	if err != nil {
		SugaredLogger.Error(err)
	}
	err = toml.Unmarshal(file, &OWC)

	if err != nil {
		panic(err)
	}
}