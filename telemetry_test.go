package kontaktparser

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildField(t *testing.T, pid TelemetryPID, dataHex string) KontaktTelemetryValue {
	data, err := hex.DecodeString(dataHex)
	if err != nil {
		t.FailNow()
	}
	return KontaktTelemetryValue{
		PID:   pid,
		Value: data,
	}
}

func TestSystemHealthField(t *testing.T) {
	tlm := buildField(t, SystemHealth, "002F685964")
	field := SystemHealthParser{}
	field.Parse(tlm)

	assert.Equal(t, uint32(1500000000), field.UnixTimestamp)
	assert.Equal(t, uint8(100), field.BatteryLevel)
}

func TestAccelerometerField(t *testing.T) {
	tlm := buildField(t, Accelerometer, "201020306400C800")
	field := AccelerometerFieldParser{}
	field.Parse(tlm)

	assert.Equal(t, uint8(32), field.Sensitivity)
	assert.Equal(t, int8(16), field.X)
	assert.Equal(t, int8(32), field.Y)
	assert.Equal(t, int8(48), field.Z)
	assert.Equal(t, uint16(100), field.SecondsSinceDoubleTap)
	assert.Equal(t, uint16(200), field.SecondsSinceThreshold)
}

func TestSensorsField(t *testing.T) {
	tlm := buildField(t, Sensors, "6410")
	field := SensorsFieldParser{}
	field.Parse(tlm)

	assert.Equal(t, uint8(100), field.LightLevel)
	assert.Equal(t, int8(16), field.Temperature)
}
