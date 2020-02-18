package kontaktparser

import (
	"encoding/binary"
	"errors"
)

var (
	ErrInvalidTelemetryPID = errors.New("invalid telemetry field pid")
)

type FieldParser interface {
	Parse(value KontaktTelemetryValue) error
}

func assertions(value KontaktTelemetryValue, pid TelemetryPID, length int) error {
	if value.PID != pid || len(value.Value) != length {
		return ErrInvalidTelemetryPID
	}
	return nil
}

type SystemHealthParser struct {
	UnixTimestamp uint32
	BatteryLevel  uint8
}

func (p *SystemHealthParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, SystemHealth, 5); err != nil {
		return err
	}
	p.UnixTimestamp = binary.LittleEndian.Uint32(value.Value[:4])
	p.BatteryLevel = uint8(value.Value[4])
	return nil
}

type AccelerometerFieldParser struct {
	Sensitivity           uint8
	X                     int8
	Y                     int8
	Z                     int8
	SecondsSinceDoubleTap uint16
	SecondsSinceThreshold uint16
}

func (p *AccelerometerFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Accelerometer, 8); err != nil {
		return err
	}
	p.Sensitivity = uint8(value.Value[0])
	p.X = int8(value.Value[1])
	p.Y = int8(value.Value[2])
	p.Z = int8(value.Value[3])
	p.SecondsSinceDoubleTap = binary.LittleEndian.Uint16(value.Value[4:6])
	p.SecondsSinceThreshold = binary.LittleEndian.Uint16(value.Value[6:8])
	return nil
}

type SensorsFieldParser struct {
	LightLevel  uint8
	Temperature int8
}

func (p *SensorsFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Sensors, 2); err != nil {
		return err
	}
	p.LightLevel = uint8(value.Value[0])
	p.Temperature = int8(value.Value[1])
	return nil
}

type AccelerationFieldParser struct {
	Sensitivity uint8
	X           int8
	Y           int8
	Z           int8
}

func (p *AccelerationFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Acceleration, 4); err != nil {
		return err
	}
	p.Sensitivity = uint8(value.Value[0])
	p.X = int8(value.Value[1])
	p.Y = int8(value.Value[2])
	p.Z = int8(value.Value[3])
	return nil
}

type MovementFieldParser struct {
	SecondsSinceThreshold uint16
}

func (p *MovementFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Movement, 2); err != nil {
		return err
	}
	p.SecondsSinceThreshold = binary.LittleEndian.Uint16(value.Value)
	return nil
}
