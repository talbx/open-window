package broker

import (
	"encoding/json"
	"errors"
	"github.com/talbx/openwindow/pkg/service"
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

var randomTopic string

func (m *messageMock) Topic() string {
	m.Called()
	return randomTopic
}

type clientMock struct {
	mock.Mock
	mqtt.Client
}

type clientMock2 struct {
	mock.Mock
	mqtt.Client
}

type tokenMock struct {
	mqtt.Token
	mock.Mock
}

type tokenMock2 struct {
	mqtt.Token
	mock.Mock
}

func (t *tokenMock) Wait() bool {
	t.Called()
	return false
}

func (t *tokenMock) Error() error {
	t.Called()
	return nil
}

func (t *tokenMock2) Wait() bool {
	t.Called()
	return true
}

func (t *tokenMock2) Error() error {
	t.Called()
	return errors.New("oh no")
}

func (c *clientMock) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	c.Called(topic)
	t := new(tokenMock)
	t.On("Wait")
	t.On("Error")
	return t
}

func (c *clientMock2) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	c.Called(topic)
	t := new(tokenMock2)
	t.On("Wait")
	t.On("Error")
	return t
}

func Test_OnConnect(t *testing.T) {
	// given
	c := new(changeMock)
	c.On("HandleChange")
	b := Broker{c, new(exitMock)}
	model.CreateSugaredLogger()
	model.OWC.Devices = []model.Device{
		{
			Topic: "any",
		},
	}
	cl := new(clientMock)
	cl.On("Subscribe", "any")

	// when
	b.OnConnect(cl)

	// then
	cl.AssertCalled(t, "Subscribe", "any")

}

func Test_createMqttOpts(t *testing.T) {
	// given
	c := new(changeMock)
	c.On("HandleChange")
	b := Broker{c, new(exitMock)}
	model.CreateSugaredLogger()
	model.OWC.MqttConfig.Host = "mqtt://1234"
	model.OWC.MqttConfig.ClientId = "my-client"

	// when
	opts := b.createMqttOpts()

	// then
	assert.NotEqual(t, "my-client", opts.ClientID)
	assert.Contains(t, opts.ClientID, "my-client")
	assert.Equal(t, "1234", opts.Servers[0].Host)
}

func Test_panic(t *testing.T) {
	c := new(changeMock)
	c.On("HandleChange")
	b := Broker{c, new(exitMock)}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// given
	model.CreateSugaredLogger()
	model.OWC.Devices = []model.Device{
		{
			Topic: "any",
		},
	}
	cl := new(clientMock2)
	c.On("Subscribe", "any")

	// when
	b.OnConnect(cl)

	// then
}

func Test_handleMessage(t *testing.T) {
	randomTopic = "testTopic"
	c := new(changeMock)
	c.On("HandleChange")
	b := Broker{c, new(exitMock)}
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

	b.HandleMessage(client, msg)

	c.AssertCalled(t, "HandleChange")
	msg.AssertExpectations(t)
	msg.AssertCalled(t, "Payload")
	msg.AssertCalled(t, "Topic")

}

type changeI interface {
	HandleChange(h model.TuyaHumidity)
}

type changeMock struct {
	service.ChangeService
	mock.Mock
}

func (c *changeMock) HandleChange(h model.TuyaHumidity) {
	c.Called()
}

func Test_handleMessage_translateErr(t *testing.T) {
	c := new(changeMock)
	c.On("HandleChange")
	b := Broker{c, new(exitMock)}
	randomTopic = "somethingNonExistent"
	change := new(changeMock)
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

	b.HandleMessage(client, msg)

	change.AssertNotCalled(t, "HandleChange")
	msg.AssertExpectations(t)
	msg.AssertCalled(t, "Payload")
	msg.AssertCalled(t, "Topic")

}

func Test_handleMessage_unmarshalErr(t *testing.T) {
	c := new(changeMock)
	c.On("HandleChange")
	b := Broker{c, new(exitMock)}
	randomTopic = "somethingNonExistent"
	change := new(changeMock)
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

	b.HandleMessage(client, msg)

	change.AssertNotCalled(t, "HandleChange")
	msg.AssertExpectations(t)
	msg.AssertCalled(t, "Payload")
	msg.AssertCalled(t, "Topic")

}

type exitMock struct {
	mock.Mock
}

func (e *exitMock) Exit(err error) {
	e.Called()
}

func Test_Connect(t *testing.T) {
	model.CreateSugaredLogger()
	exiter := new(exitMock)
	b := Broker{new(changeMock), exiter}
	tm1 := new(tokenMock)
	tm1.On("Wait")
	tm1.On("Error")

	// when
	b.Connect(tm1)

	// then
	tm1.AssertCalled(t, "Wait")
	tm1.AssertNotCalled(t, "Error")
	exiter.AssertNotCalled(t, "Exit")
}

func Test_Connect_fail(t *testing.T) {

	model.CreateSugaredLogger()
	exiter := new(exitMock)
	exiter.On("Exit")
	b := Broker{new(changeMock), exiter}
	tm2 := new(tokenMock2)
	tm2.On("Wait")
	tm2.On("Error")

	// when
	b.Connect(tm2)

	// then
	tm2.AssertCalled(t, "Wait")
	tm2.AssertCalled(t, "Error")

	exiter.AssertCalled(t, "Exit")

}
