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

func assertions(value KontaktTelemetryValue, pid TelemetryPid, len int) error {
	if value.PID != pid || len(value.Value) != len {
		return ErrInvalidTelemetryPID
	}
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
