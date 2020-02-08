package kontaktparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/google/uuid"
)

const (
	// Unknown - packet of not known origin, that doesn't match any of the supported ones
	Unknown = iota
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
	ErrInvalidPreamble = errors.New("Invalid preamble")
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
	DetectedType int
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
			if ibeacon, err := parseIBeacon(section); err == nil {
				p.Parsed = ibeacon
			} else if err == ErrInvalidPreamble {
				continue
			} else {
				return err
			}
		case serviceDataDataType:

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

func parseIBeacon(section []byte) (*IBeaconAdvertisement, error) {
	if !bytes.Equal(section[0:4], ibeaconManufacturerConstData) {
		return nil, ErrInvalidPreamble
	}
	proximity, err := uuid.FromBytes(section[4:20])
	if err != nil {
		return nil, err
	}
	major := section[20:22]
	minor := section[22:24]
	rssi := section[24]
	return &IBeaconAdvertisement{
		CalibratedRssi: int8(rssi),
		ProximityUUID:  proximity,
		Major:          binary.LittleEndian.Uint16(major),
		Minor:          binary.LittleEndian.Uint16(minor),
	}, nil
}
