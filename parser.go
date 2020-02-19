package kontaktparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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
	ErrInvalidURL                      = errors.New("invalid eddystone url")
)

var (
	flagsDataType                byte = 0x01
	completeNameType             byte = 0x09
	txPowerType                  byte = 0x0A
	serviceDataDataType          byte = 0x16
	manufacturerDataType         byte = 0xFF
	ibeaconManufacturerConstData      = []byte{0x4C, 0x00, 0x02, 0x15}
	kontaktUUID                       = []byte{0x6A, 0xFE}
	eddystoneUUID                     = []byte{0xAA, 0xFE}
	kontaktScanResponseUUID           = []byte{0x0D, 0xD0}
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

func (p *Parser) ParseScanResponse() error {
	scanResponse := KontaktIOScanResponse{}
	for p.buf.Len() > 0 {
		typ, section, err := p.nextSection()
		if err != nil {
			return err
		}
		switch typ {
		case completeNameType:
			scanResponse.Name = string(section)
			scanResponse.HasName = true
		case txPowerType:
			scanResponse.TxPower = int8(section[0])
			scanResponse.HasTxPower = true
		case serviceDataDataType:
			if !bytes.Equal(section[0:2], kontaktScanResponseUUID) || len(section) != 9 {
				continue
			}
			scanResponse.UniqueID = string(section[2:6])
			scanResponse.Firmware = fmt.Sprintf("%v.%v", section[6], section[7])
			scanResponse.BatteryLevel = uint8(section[8])
			scanResponse.HasIdentifier = true
		}
	}
	if scanResponse.HasName || scanResponse.HasTxPower || scanResponse.HasIdentifier {
		p.Parsed = &scanResponse
		p.DetectedType = KontaktScanResponse
	}
	return nil
}

func (p *Parser) ParseAdvertisement() error {
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
			if len(section) < 2 {
				return io.EOF
			}
			uuid := section[0:2]
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
	p.Parsed = &KontaktPlainAdvertisement{
		DeviceModel:   uint8(section[3]),
		FirmwareMajor: uint8(section[4]),
		FirmwareMinor: uint8(section[5]),
		BatteryLevel:  uint8(section[6]),
		TxPower:       int8(section[7]),
		UniqueID:      string(section[8:]),
	}
	p.DetectedType = KontaktPlain
	return nil
}

func (p *Parser) parseKontaktShuffled(section []byte) error {
	if len(section) != 23 {
		return nil
	}
	p.Parsed = &KontaktShuffledAdvertisement{
		DeviceModel:         uint8(section[3]),
		FirmwareMajor:       uint8(section[4]),
		FirmwareMinor:       uint8(section[5]),
		BatteryLevel:        uint8(section[6]),
		TxPower:             int8(section[7]),
		EddystoneNamespace:  section[8:18],
		EddystoneInstanceID: section[18:24],
	}
	p.DetectedType = KontaktShuffled
	return nil
}

func (p *Parser) parseKontaktTelemetry(section []byte) error {
	fields := make([]KontaktTelemetryValue, 0)
	buf := bytes.NewBuffer(section[3:])
	for buf.Len() != 0 {
		len, err := buf.ReadByte()
		if err != nil || buf.Len() < int(len) {
			return io.EOF
		}
		pid, err := buf.ReadByte()
		if err != nil {
			return err
		}
		value := make([]byte, len-1)
		if n, err := buf.Read(value); err != nil {
			return err
		} else if n != int(len)-1 {
			return io.EOF
		}
		fields = append(fields, KontaktTelemetryValue{
			PID:   TelemetryPID(pid),
			Value: value,
		})
	}
	p.Parsed = &KontaktTelemetryAdvertisement{Fields: fields}
	p.DetectedType = KontaktTelemetry
	return nil
}

func (p *Parser) parseKontaktLocation(section []byte) error {
	if len(section) < 8 {
		return nil
	}

	txPower := section[3]
	bleChannel := section[4]
	model := section[5]
	flags := section[6]
	uniqueID := string(section[7:])

	p.Parsed = &KontaktLocationAdvertisement{
		TxPower:     int8(txPower),
		BleChannel:  uint8(bleChannel),
		DeviceModel: uint8(model),
		Flags:       uint8(flags),
		UniqueID:    uniqueID,
	}
	p.DetectedType = KontaktLocation
	return nil
}

func (p *Parser) parseEddystone(section []byte) error {
	if len(section) < 3 {
		return io.EOF
	}
	var err error
	switch section[2] {
	case 0x00:
		err = p.parseEddystoneUID(section)
	case 0x10:
		err = p.parseEddystoneURL(section)
	case 0x20:
		err = p.parseEddystoneTLM(section)
	case 0x30:
		err = p.parseEddystoneEID(section)
	}
	return err
}

func (p *Parser) parseEddystoneUID(section []byte) error {
	if len(section) != 22 {
		return io.EOF
	}
	p.Parsed = &EddystoneUIDPacket{
		TxPower0M:  int8(section[3]),
		Namespace:  section[4:14],
		InstanceId: section[14:20],
	}
	p.DetectedType = EddystoneUID
	return nil
}

var eddystoneUrlPrefixes = map[byte][]byte{
	0x00: []byte("http://www."),
	0x01: []byte("https://www."),
	0x02: []byte("http://"),
	0x03: []byte("https://"),
}

var eddystoneUrlReplacements = map[byte][]byte{
	0x00: []byte(".com/"),
	0x01: []byte(".org/"),
	0x02: []byte(".edu/"),
	0x03: []byte(".net/"),
	0x04: []byte(".info/"),
	0x05: []byte(".biz/"),
	0x06: []byte(".gov/"),
	0x07: []byte(".com"),
	0x08: []byte(".org"),
	0x09: []byte(".edu"),
	0x0A: []byte(".net"),
	0x0B: []byte(".info"),
	0x0C: []byte(".biz"),
	0x0D: []byte(".gov"),
}

func (p *Parser) parseEddystoneURL(section []byte) error {
	if len(section) < 6 {
		return io.EOF
	}
	txPower := int8(section[3])
	url := make([]byte, 0)
	if prefix, ok := eddystoneUrlPrefixes[section[4]]; ok {
		url = append(url, prefix...)
	} else {
		return ErrInvalidURL
	}
	for i := 5; i < len(section); i++ {
		b := section[i]
		if b >= 0x20 && b <= 0x7F {
			url = append(url, b)
			continue
		}
		if replacement, ok := eddystoneUrlReplacements[b]; ok {
			url = append(url, replacement...)
		} else {
			return ErrInvalidURL
		}
	}
	p.Parsed = &EddystoneURLPacket{
		TxPower0M: txPower,
		URL:       string(url),
	}
	p.DetectedType = EddystoneURL
	return nil
}

func (p *Parser) parseEddystoneTLM(section []byte) error {
	return ErrNotImplemented
}

func (p *Parser) parseEddystoneEID(section []byte) error {
	return ErrNotImplemented
}
