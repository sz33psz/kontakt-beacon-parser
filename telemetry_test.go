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
	field := SystemHealthFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint32(1500000000), field.UnixTimestamp)
	assert.Equal(t, uint8(100), field.BatteryLevel)
}

func TestSystemHealthWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := SystemHealthFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestAccelerometerField(t *testing.T) {
	tlm := buildField(t, Accelerometer, "201020306400C800")
	field := AccelerometerFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(32), field.Sensitivity)
	assert.Equal(t, int8(16), field.X)
	assert.Equal(t, int8(32), field.Y)
	assert.Equal(t, int8(48), field.Z)
	assert.Equal(t, uint16(100), field.SecondsSinceDoubleTap)
	assert.Equal(t, uint16(200), field.SecondsSinceThreshold)
}

func TestAccelerometerWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := AccelerometerFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestSensorsField(t *testing.T) {
	tlm := buildField(t, Sensors, "6410")
	field := SensorsFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(100), field.LightLevel)
	assert.Equal(t, int8(16), field.Temperature)
}

func TestSensorsWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := SensorsFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestAccelerationField(t *testing.T) {
	tlm := buildField(t, Acceleration, "20102030")
	field := AccelerationFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(32), field.Sensitivity)
	assert.Equal(t, int8(16), field.X)
	assert.Equal(t, int8(32), field.Y)
	assert.Equal(t, int8(48), field.Z)
}

func TestAccelerationWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := AccelerationFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestMovementField(t *testing.T) {
	tlm := buildField(t, Movement, "6401")
	field := MovementFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint16(356), field.SecondsSinceThreshold)
}

func TestMovementWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := MovementFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestDoubleTapField(t *testing.T) {
	tlm := buildField(t, DoubleTap, "6401")
	field := DoubleTapFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint16(356), field.SecondsSinceDoubleTap)
}

func TestDoubleTapWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := DoubleTapFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestLightLevelField(t *testing.T) {
	tlm := buildField(t, LightLevel, "14")
	field := LightLevelFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(20), field.LightLevel)
}

func TestLightLevelWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := LightLevelFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestTemperature8BitField(t *testing.T) {
	tlm := buildField(t, Temperature8Bit, "FE")
	field := Temperature8BitFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, int8(-2), field.Temperature)
}

func TestTemperature8BitWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := Temperature8BitFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestTemperature16BitField(t *testing.T) {
	tlm := buildField(t, Temperature16Bit, "FD80")
	field := Temperature16BitFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, float32(-2.5), field.Temperature)
}

func TestTemperature16BitWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := Temperature16BitFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestBatteryField(t *testing.T) {
	tlm := buildField(t, BatteryLevel, "40")
	field := BatteryFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(64), field.BatteryLevel)
}

func TestBatteryWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := BatteryFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestClickField(t *testing.T) {
	tlm := buildField(t, Click, "6401")
	field := ClickFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint16(356), field.SecondsSinceClick)
}

func TestClickWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := ClickFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestClickInfoField(t *testing.T) {
	tlm := buildField(t, ClickInfo, "406401")
	field := ClickInfoFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(64), field.ClickID)
	assert.Equal(t, uint16(356), field.SecondsSinceClick)
}

func TestClickInfoWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := ClickInfoFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestUTCTimeField(t *testing.T) {
	tlm := buildField(t, UTCTime, "002F6859")
	field := UTCTimeFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint32(1500000000), field.UTCTime)
}

func TestUTCTimeWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := UTCTimeFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestHumidityField(t *testing.T) {
	tlm := buildField(t, Humidity, "24")
	field := HumidityFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(36), field.Humidity)
}

func TestHumidityWrongPID(t *testing.T) {
	tlm := buildField(t, Temperature8Bit, "00")
	field := HumidityFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}

func TestMovementInfoField(t *testing.T) {
	tlm := buildField(t, MovementInfo, "246401")
	field := MovementInfoFieldParser{}
	assert.Nil(t, field.Parse(tlm))

	assert.Equal(t, uint8(36), field.Counter)
	assert.Equal(t, uint16(356), field.SecondsSinceThreshold)
}

func TestMovementInfoWrongPID(t *testing.T) {
	tlm := buildField(t, Humidity, "00")
	field := MovementInfoFieldParser{}
	assert.Equal(t, ErrInvalidTelemetryPID, field.Parse(tlm))
}
