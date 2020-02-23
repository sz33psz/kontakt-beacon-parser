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

func TestParseTelemetryFrameOneField(t *testing.T) {
	bytes, err := hex.DecodeString("07166AFE03020A64")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, KontaktTelemetry, parser.DetectedType)

	if adv, ok := parser.Parsed.(*KontaktTelemetryAdvertisement); !ok {
		t.Errorf("Parsing of kontakt telemetry should result in KontaktTelemetryAdvertisement")
	} else {
		assert.Equal(t, 1, len(adv.Fields))
		assert.Equal(t, LightLevel, adv.Fields[0].PID)
		assert.Equal(t, []byte{0x64}, adv.Fields[0].Value)
	}
}

func TestParseTelemetryFrameManyFields(t *testing.T) {
	bytes, err := hex.DecodeString("0C166AFE03020A640411065BA0")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, KontaktTelemetry, parser.DetectedType)

	if adv, ok := parser.Parsed.(*KontaktTelemetryAdvertisement); !ok {
		t.Errorf("Parsing of kontakt telemetry should result in KontaktTelemetryAdvertisement")
	} else {
		assert.Equal(t, 2, len(adv.Fields))
		assert.Equal(t, LightLevel, adv.Fields[0].PID)
		assert.Equal(t, []byte{0x64}, adv.Fields[0].Value)
		assert.Equal(t, ClickInfo, adv.Fields[1].PID)
		assert.Equal(t, []byte{0x06, 0x5B, 0xA0}, adv.Fields[1].Value)
	}
}

func TestParseTelemetryFrameEmpty(t *testing.T) {
	bytes, err := hex.DecodeString("04166AFE03")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, KontaktTelemetry, parser.DetectedType)

	if adv, ok := parser.Parsed.(*KontaktTelemetryAdvertisement); !ok {
		t.Errorf("Parsing of kontakt telemetry should result in KontaktTelemetryAdvertisement")
	} else {
		assert.Equal(t, 0, len(adv.Fields))
	}
}

func TestParseEddystoneUID(t *testing.T) {
	bytes, err := hex.DecodeString("1716AAFE0004010203040506070809000102030405060000")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, EddystoneUID, parser.DetectedType)

	if adv, ok := parser.Parsed.(*EddystoneUIDPacket); !ok {
		t.Errorf("Parsing of eddystone uid should result in EddystoneUIDPacket")
	} else {
		assert.Equal(t, int8(4), adv.TxPower0M)
		assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00}, adv.Namespace)
		assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, adv.InstanceId)
	}
}

func TestEddystoneURL(t *testing.T) {
	bytes, err := hex.DecodeString("0B16AAFE100403746573740C")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, EddystoneURL, parser.DetectedType)

	if adv, ok := parser.Parsed.(*EddystoneURLPacket); !ok {
		t.Errorf("Parsing of eddystone url should result in EddystoneURLPacket")
	} else {
		assert.Equal(t, int8(4), adv.TxPower0M)
		assert.Equal(t, "https://test.biz", adv.URL)
	}
}

func TestEddystoneURL2(t *testing.T) {
	bytes, err := hex.DecodeString("0F16AAFE100400746573740674657374")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, EddystoneURL, parser.DetectedType)

	if adv, ok := parser.Parsed.(*EddystoneURLPacket); !ok {
		t.Errorf("Parsing of eddystone url should result in EddystoneURLPacket")
	} else {
		assert.Equal(t, int8(4), adv.TxPower0M)
		assert.Equal(t, "http://www.test.gov/test", adv.URL)
	}
}

func TestEddystoneTLM(t *testing.T) {
	bytes, err := hex.DecodeString("1116AAFE2000018005400000010000010000")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, EddystoneTLM, parser.DetectedType)

	if adv, ok := parser.Parsed.(*EddystonePlainTLMPacket); !ok {
		t.Errorf("Parsing of eddystone tlm should result in EddystonePlainTLMPacket")
	} else {
		assert.Equal(t, uint16(384), adv.BatteryVoltage)
		assert.Equal(t, float64(5.25), adv.Temperature)
		assert.Equal(t, uint32(256), adv.AdvertisementCount)
		assert.Equal(t, float64(6553.6), adv.TimeSincePowerOn)
	}
}

func TestEddystoneETLM(t *testing.T) {
	bytes, err := hex.DecodeString("1516AAFE20010102030405060708090A0B0C01021112")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, EddystoneETLM, parser.DetectedType)

	if adv, ok := parser.Parsed.(*EddystoneEncryptedTLMPacket); !ok {
		t.Errorf("Parsing of eddystone tlm should result in EddystoneEncryptedTLMPacket")
	} else {
		assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C}, adv.Telemetry)
		assert.Equal(t, []byte{0x01, 0x02}, adv.Salt)
		assert.Equal(t, []byte{0x11, 0x12}, adv.MIC)
	}
}

func TestEddystoneEID(t *testing.T) {
	bytes, err := hex.DecodeString("0D16AAFE30045152535455565758")
	assert.Nil(t, err)

	parser := New(bytes)
	assert.Nil(t, parser.ParseAdvertisement())
	assert.Equal(t, EddystoneEID, parser.DetectedType)

	if adv, ok := parser.Parsed.(*EddystoneEIDPacket); !ok {
		t.Errorf("Parsing of eddystone eid should result in EddystoneEIDPacket")
	} else {
		assert.Equal(t, int8(4), adv.TxPower0M)
		assert.Equal(t, []byte{0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58}, adv.EID)
	}
}
