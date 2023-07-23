package model

import (
	"github.com/k0kubun/pp"
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
	Topic string `yaml:"topic"`
	Room  string `yaml:"room"`
}

type PushoverConfig struct {
	ApiToken  string `yaml:"apiToken"`
	UserToken string `yaml:"userToken"`
}

type MqttConfig struct {
	Host     string `yaml:"host"`
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

	SugaredLogger.Infof("successfully loaded OWC config")
	logConfig()
	if err != nil {
		SugaredLogger.Errorf("there was an error parsing the config.toml", err)
		os.Exit(1)
	}
}

func logConfig() {
	pp.Printf("mqtt-conf: %+v\n", OWC.MqttConfig)
	pp.Printf("openwindow-conf: %+v\n", OWC.OpenWindowConfig)

	oapi := obfucscate(OWC.ApiToken)
	ouse := obfucscate(OWC.UserToken)
	pp.Printf("pushover-conf: %+v\n", PushoverConfig{oapi, ouse})
}

func obfucscate(s string) string {
	obfuscate := ""
	for i := 0; i < len(s); i++ {
		obfuscate += "*"
	}
	return obfuscate
}
