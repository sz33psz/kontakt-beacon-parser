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

type SystemHealthFieldParser struct {
	UnixTimestamp uint32
	BatteryLevel  uint8
}

func (p *SystemHealthFieldParser) Parse(value KontaktTelemetryValue) error {
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

type DoubleTapFieldParser struct {
	SecondsSinceDoubleTap uint16
}

func (p *DoubleTapFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, DoubleTap, 2); err != nil {
		return err
	}
	p.SecondsSinceDoubleTap = binary.LittleEndian.Uint16(value.Value)
	return nil
}

type LightLevelFieldParser struct {
	LightLevel uint8
}

func (p *LightLevelFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, LightLevel, 1); err != nil {
		return err
	}
	p.LightLevel = uint8(value.Value[0])
	return nil
}

type Temperature8BitFieldParser struct {
	Temperature int8
}

func (p *Temperature8BitFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Temperature8Bit, 1); err != nil {
		return err
	}
	p.Temperature = int8(value.Value[0])
	return nil
}

type Temperature16BitFieldParser struct {
	Temperature float32
}

func (p *Temperature16BitFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Temperature16Bit, 2); err != nil {
		return err
	}
	p.Temperature = float32((int16(value.Value[0])<<8)+int16(value.Value[1])) / 256

	return nil
}

type BatteryFieldParser struct {
	BatteryLevel uint8
}

func (p *BatteryFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, BatteryLevel, 1); err != nil {
		return err
	}
	p.BatteryLevel = uint8(value.Value[0])
	return nil
}

type ClickFieldParser struct {
	SecondsSinceClick uint16
}

func (p *ClickFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Click, 2); err != nil {
		return err
	}
	p.SecondsSinceClick = binary.LittleEndian.Uint16(value.Value)
	return nil
}

type ClickInfoFieldParser struct {
	ClickID           uint8
	SecondsSinceClick uint16
}

func (p *ClickInfoFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, ClickInfo, 3); err != nil {
		return err
	}
	p.ClickID = uint8(value.Value[0])
	p.SecondsSinceClick = binary.LittleEndian.Uint16(value.Value[1:])
	return nil
}

type UTCTimeFieldParser struct {
	UTCTime uint32
}

func (p *UTCTimeFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, UTCTime, 4); err != nil {
		return err
	}
	p.UTCTime = binary.LittleEndian.Uint32(value.Value)
	return nil
}

type HumidityFieldParser struct {
	Humidity uint8
}

func (p *HumidityFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, Humidity, 1); err != nil {
		return err
	}
	p.Humidity = uint8(value.Value[0])
	return nil
}

type MovementInfoFieldParser struct {
	Counter               uint8
	SecondsSinceThreshold uint16
}

func (p *MovementInfoFieldParser) Parse(value KontaktTelemetryValue) error {
	if err := assertions(value, MovementInfo, 3); err != nil {
		return err
	}
	p.Counter = uint8(value.Value[0])
	p.SecondsSinceThreshold = binary.LittleEndian.Uint16(value.Value[1:])
	return nil
}
