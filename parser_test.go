package kontaktparser

import (
	"encoding/hex"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseIBeacon(t *testing.T) {
	bytes, err := hex.DecodeString("1AFF4C000215F7826DA64FA24E988024BC5B71E0893E01020304B3")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, parser.DetectedType, IBeacon)
	if adv, ok := parser.Parsed.(*IBeaconAdvertisement); !ok {
		t.Errorf("Parsing of iBeacon should result in IBeaconAdvertisement")
	} else {
		assert.Equal(t, uuid.MustParse("F7826DA6-4FA2-4E98-8024-BC5B71E0893E"), adv.ProximityUUID)
		assert.Equal(t, uint16(513), adv.Major)
		assert.Equal(t, uint16(1027), adv.Minor)
		assert.Equal(t, int8(-77), adv.CalibratedRssi)
	}
}
func TestParseIBeaconInvalidPreamble(t *testing.T) {
	bytes, err := hex.DecodeString("1AFFFFFFFFFFF7826DA64FA24E988024BC5B71E0893E01020304B3")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, parser.DetectedType, Unknown)
}
func TestParseIBeaconSmallerFrameLength(t *testing.T) {
	bytes, err := hex.DecodeString("19FFFFFFFFFFF7826DA64FA24E988024BC5B71E0893E01020304")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, parser.DetectedType, Unknown)
}

func TestParseTooShortIBeacon(t *testing.T) {
	bytes, err := hex.DecodeString("1AFF4C000215F7826DA64FA24E988024BC5B71E0893E01020304")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Equal(t, io.EOF, parser.ParseAdvertisement())
	assert.Equal(t, Unknown, parser.DetectedType)
}

func TestParseKontaktInvalidType(t *testing.T) {
	bytes, err := hex.DecodeString("0F166AFEFF06010F6404616263646566")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Equal(t, ErrInvalidKontaktPayloadIdentifier, parser.ParseAdvertisement())
	assert.Equal(t, Unknown, parser.DetectedType)
}

func TestParseKontaktTooShort(t *testing.T) {
	bytes, err := hex.DecodeString("02166A")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Equal(t, io.EOF, parser.ParseAdvertisement())
	assert.Equal(t, Unknown, parser.DetectedType)
}

func TestParseBlankAdvertisement(t *testing.T) {
	// It can cause infinite loop when reading 0xFF as -1
	bytes, err := hex.DecodeString("FFFFFFFFFFFFFFFFFFFFFF")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Equal(t, io.EOF, parser.ParseAdvertisement())
	assert.Equal(t, Unknown, parser.DetectedType)
}

func TestParseKontaktPlain(t *testing.T) {
	bytes, err := hex.DecodeString("0F166AFE0206010F6404616263646566")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, KontaktPlain, parser.DetectedType)
	if adv, ok := parser.Parsed.(*KontaktPlainAdvertisement); !ok {
		t.Errorf("Parsing of iBeacon should result in KontaktPlainAdvertisement")
	} else {
		assert.Equal(t, uint8(100), adv.BatteryLevel)
		assert.Equal(t, uint8(1), adv.FirmwareMajor)
		assert.Equal(t, uint8(15), adv.FirmwareMinor)
		assert.Equal(t, uint8(6), adv.DeviceModel)
		assert.Equal(t, "abcdef", adv.UniqueID)
	}
}

func TestParseKontaktPlainTooShort(t *testing.T) {
	bytes, err := hex.DecodeString("09166AFE0206010F6404")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, Unknown, parser.DetectedType)
}

func TestParseScanResponse(t *testing.T) {
	bytes, err := hex.DecodeString("080961626364656667020A040A160DD061626364040264")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseScanResponse())
	assert.Equal(t, KontaktScanResponse, parser.DetectedType)
	if sr, ok := parser.Parsed.(*KontaktIOScanResponse); !ok {
		t.Errorf("Parsing of kontakt scan response should result in KontaktIOScanResponse")
	} else {
		assert.True(t, sr.HasName)
		assert.Equal(t, "abcdefg", sr.Name)
		assert.True(t, sr.HasTxPower)
		assert.Equal(t, int8(4), sr.TxPower)
		assert.True(t, sr.HasIdentifier)
		assert.Equal(t, "abcd", sr.UniqueID)
		assert.Equal(t, "4.2", sr.Firmware)
		assert.Equal(t, uint8(100), sr.BatteryLevel)
	}
}

func TestParseScanResponseTooShortSection(t *testing.T) {
	bytes, err := hex.DecodeString("080961626364656667020A0409160DD0616263640402")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseScanResponse())
	assert.Equal(t, KontaktScanResponse, parser.DetectedType)
	if sr, ok := parser.Parsed.(*KontaktIOScanResponse); !ok {
		t.Errorf("Parsing of kontakt scan response should result in KontaktIOScanResponse")
	} else {
		assert.True(t, sr.HasName)
		assert.Equal(t, "abcdefg", sr.Name)
		assert.True(t, sr.HasTxPower)
		assert.Equal(t, int8(4), sr.TxPower)
		assert.False(t, sr.HasIdentifier)
	}
}

func TestParseLocationFrame(t *testing.T) {
	bytes, err := hex.DecodeString("0E166AFE05F4250A01414243444546")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, KontaktLocation, parser.DetectedType)

	if adv, ok := parser.Parsed.(*KontaktLocationAdvertisement); !ok {
		t.Errorf("Parsing of kontakt location should result in KontaktLocationAdvertisement")
	} else {
		assert.Equal(t, int8(-12), adv.TxPower)
		assert.Equal(t, uint8(37), adv.BleChannel)
		assert.Equal(t, uint8(10), adv.DeviceModel)
		assert.Equal(t, uint8(1), adv.Flags)
		assert.Equal(t, "ABCDEF", adv.UniqueID)
	}
}

func TestParseLocationFrameTooShort(t *testing.T) {
	bytes, err := hex.DecodeString("08166AFE05F4250A01")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, Unknown, parser.DetectedType)
}
