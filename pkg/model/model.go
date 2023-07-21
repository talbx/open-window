package model

import (
	"os"

	"gopkg.in/yaml.v3"
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
	Room  string
}

type PushoverConfig struct {
	ApiToken  string
	UserToken string
}

type MqttConfig struct {
	Host     string
	ClientId string `yaml:"clientId"`
}

type GlobalConfig struct {
	PushoverConfig   `yaml:"pushover"`
	OpenWindowConfig `yaml:"openwindow"`
	MqttConfig       `yaml:"mqtt"`
}

type OpenWindowConfig struct {
	Devices  []Device
	Interval int
}

var OWC GlobalConfig

func GetGlobalConfig() *GlobalConfig {
	return &OWC
}

func CreateOpenWindowConfig() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		SugaredLogger.Error(err)
		SugaredLogger.Error("No config.yaml provided, will exit now!")
		os.Exit(1)
	}
	err = yaml.Unmarshal(file, &OWC)

	if err != nil {
		SugaredLogger.Errorf("there was an error parsing the config.toml", err)
		os.Exit(1)
	}
}
