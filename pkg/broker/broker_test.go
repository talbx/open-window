package broker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/talbx/openwindow/pkg/model"
)

func Test_translateName(t *testing.T) {
	model.OWC.Devices = []model.Device{
		model.Device{
			Topic: "testTopic",
			Room: "testRoom",
		},
		model.Device{
			Topic: "someOtherTopic",
			Room: "otherRoom",
		},
	}

	room, err := translateName("testTopic")
	
}
