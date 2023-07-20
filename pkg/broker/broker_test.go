package broker

import (
	"encoding/json"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/talbx/openwindow/pkg/model"
)

func Test_translateName_succ(t *testing.T) {
	model.OWC.Devices = []model.Device{
		{
			Topic: "testTopic",
			Room:  "testRoom",
		},
		{
			Topic: "someOtherTopic",
			Room:  "otherRoom",
		},
	}

	room, err := Translate("testTopic")
	assert.Equal(t, "testRoom", room)
	assert.Nil(t, err)

	room, err = Translate("someOtherTopic")
	assert.Equal(t, "otherRoom", room)
	assert.Nil(t, err)

}

func Test_translateName_err(t *testing.T) {
	model.OWC.Devices = []model.Device{
		{
			Topic: "testTopic",
			Room:  "testRoom",
		},
	}

	room, err := Translate("somethingNonExistent")

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Empty(t, room)
}

type messageMock struct {
	mqtt.Message
	mock.Mock
}

func (m *messageMock) Payload() []byte {
	r, _ := json.Marshal(model.TuyaHumidity{})
	m.Called()
	return r
}

func (m *messageMock) Topic() string {
	m.Called()
	return "testTopic"
}

func Test_handleMessage(t *testing.T) {
	model.OWC.Devices = []model.Device{
		{
			Topic: "testTopic",
			Room:  "testRoom",
		},
	}
	model.CreateSugaredLogger()

	var client mqtt.Client
	msg := new(messageMock)
	msg.On("Payload")
	msg.On("Topic")

	handleMessage(client, msg)

	msg.AssertExpectations(t)
	msg.AssertCalled(t, "Payload")
	msg.AssertCalled(t, "Topic")

}
