package kontaktparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/google/uuid"
)

type DetectedType int

const (
	// Unknown - packet of not known origin, that doesn't match any of the supported ones
	Unknown DetectedType = iota
	// IBeacon - iBeacon packet
	IBeacon
	// EddystoneUID - Eddystone's UID packet
	EddystoneUID
	// EddystoneURL - Eddystone's URL packet
	EddystoneURL
	// EddystoneTLM - Eddystone's Telemetry packet
	EddystoneTLM
	// EddystoneEID - Eddystone's Ephemeral ID packet
	EddystoneEID
	// EddystoneETLM - Eddystone's encrypted telemetry packet
	EddystoneETLM
	// KontaktScanResponse - Kontakt.io old beacon scan response packet
	KontaktScanResponse
	// KontaktPlain - Kontakt.io Secure Profile plain advertisement packet
	KontaktPlain
	// KontaktShuffled - Kontakt.io Secure Profile shuffled advertisement packet
	KontaktShuffled
	// KontaktTelemetry - Kontakt.io Telemetry packet
	KontaktTelemetry
	// KontaktLocation - Kontakt.io Location packet
	KontaktLocation
)

var (
	ibeaconLength = 25
)

var (
	ErrInvalidPreamble                 = errors.New("invalid preamble")
	ErrInvalidKontaktPayloadIdentifier = errors.New("invalid kontakt payload identidier")
	ErrNotImplemented                  = errors.New("packet not supported yet")
	ErrInvalidLength                   = errors.New("packet has invalid length")
)

var (
	serviceDataDataType          byte = 0x16
	manufacturerDataType         byte = 0xFF
	flagsDataType                byte = 0x01
	ibeaconManufacturerConstData      = []byte{0x4C, 0x00, 0x02, 0x15}
	kontaktUUID                       = []byte{0x6A, 0xFE}
	eddystoneUUID                     = []byte{0xAA, 0xFE}
)

type Parser struct {
	buf          *bytes.Buffer
	DetectedType DetectedType
	Flags        byte
	Parsed       interface{}
}

func New(adv []byte) Parser {
	return Parser{
		buf:          bytes.NewBuffer(adv),
		DetectedType: Unknown,
	}
}

func (p *Parser) Parse() error {
	for p.buf.Len() > 0 {
		typ, section, err := p.nextSection()
		if err != nil {
			return err
		}
		switch typ {
		case flagsDataType:
			p.Flags = section[0]
		case manufacturerDataType:
			if len(section) != ibeaconLength {
				continue
			}
			if err = p.parseIBeacon(section); err == ErrInvalidPreamble {
				continue
			} else if err != nil {
				return err
			}
		case serviceDataDataType:
			if p.buf.Len() < 2 {
				return io.EOF
			}
			uuid := p.buf.Next(2)
			if bytes.Equal(uuid, kontaktUUID) {
				if err := p.parseKontaktAdv(section); err != nil {
					return err
				}
			} else if bytes.Equal(uuid, eddystoneUUID) {
				if err := p.parseEddystone(section); err != nil {
					return err
				}
			} else {
				continue
			}
		default:
			continue
		}
	}
	return nil
}

func (p *Parser) nextSection() (byte, []byte, error) {
	len, err := p.buf.ReadByte()
	if err != nil {
		if err == io.EOF {

		}
		return 0, nil, err
	}
	typ, err := p.buf.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	sectionData := make([]byte, len-1)
	n, err := p.buf.Read(sectionData)
	if err != nil {
		return typ, nil, err
	}
	if n != int(len)-1 {
		return typ, nil, io.EOF
	}
	return typ, sectionData, nil
}

func (p *Parser) parseIBeacon(section []byte) error {
	if !bytes.Equal(section[0:4], ibeaconManufacturerConstData) {
		return ErrInvalidPreamble
	}
	proximity, err := uuid.FromBytes(section[4:20])
	if err != nil {
		return err
	}
	major := section[20:22]
	minor := section[22:24]
	rssi := section[24]

	p.Parsed = &IBeaconAdvertisement{
		CalibratedRssi: int8(rssi),
		ProximityUUID:  proximity,
		Major:          binary.LittleEndian.Uint16(major),
		Minor:          binary.LittleEndian.Uint16(minor),
	}
	p.DetectedType = IBeacon
	return nil
}

func (p *Parser) parseKontaktAdv(section []byte) error {
	if len(section) < 3 {
		return io.EOF
	}
	var err error
	switch section[2] {
	case 0x01:
		err = p.parseKontaktShuffled(section)
	case 0x02:
		err = p.parseKontaktPlain(section)
	case 0x03:
		err = p.parseKontaktTelemetry(section)
	case 0x05:
		err = p.parseKontaktLocation(section)
	default:
		return ErrInvalidKontaktPayloadIdentifier
	}
	return err
}

func (p *Parser) parseKontaktPlain(section []byte) error {
	if len(section) < 9 {
		return nil
	}
	deviceModel := section[3]
	fwMajor := section[4]
	fwMinor := section[5]
	battery := section[6]
	txPower := section[7]
	uniqueID := string(section[8:])
	p.Parsed = &KontaktPlainAdvertisement{
		DeviceModel:   uint8(deviceModel),
		FirmwareMajor: uint8(fwMajor),
		FirmwareMinor: uint8(fwMinor),
		BatteryLevel:  uint8(battery),
		TxPower:       int8(txPower),
		UniqueID:      uniqueID,
	}
	p.DetectedType = KontaktPlain
	return nil
}

func (p *Parser) parseKontaktShuffled(section []byte) error {
	return ErrNotImplemented
}

func (p *Parser) parseKontaktTelemetry(section []byte) error {
	return ErrNotImplemented
}

func (p *Parser) parseKontaktLocation(section []byte) error {
	return ErrNotImplemented
}

func (p *Parser) parseEddystone(section []byte) error {
	return ErrNotImplemented
}
