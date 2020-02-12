package kontaktparser

import (
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseIBeacon(t *testing.T) {
	bytes, err := hex.DecodeString("1AFF4C000215F7826DA64FA24E988024BC5B71E0893E01020304B3")
	if err != nil {
		t.Fail()
	}

	parser := New(bytes)
	assert.Nil(t, parser.Parse())
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
