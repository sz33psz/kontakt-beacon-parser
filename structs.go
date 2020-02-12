package kontaktparser

import "github.com/google/uuid"

// IBeaconAdvertisement is a structure holding data from iBeacon advertisement
type IBeaconAdvertisement struct {
	CalibratedRssi int8
	ProximityUUID  uuid.UUID
	Major          uint16
	Minor          uint16
}

// KontaktIOScanResponse is a structure holding data from older Kontakt.io beacon's Scan Response
type KontaktIOScanResponse struct {
	Name            string
	TxPower         uint8
	Firmware        string
	BatteryLevel    uint8
	UniqueID        string
	ShuffledIBeacon IBeaconAdvertisement
}

// KontaktPlainAdvertisement is a structure holding data from Kontakt.io Secure Profile plain advertisement
// Secure profile is described here: https://developer.kontakt.io/hardware/packets/secureprofile/
type KontaktPlainAdvertisement struct {
	DeviceModel   uint8
	FirmwareMajor uint8
	FirmwareMinor uint8
	BatteryLevel  uint8
	TxPower       int8
	UniqueID      string
}

// KontaktShuffledAdvertisement is a structure holding data from Kontakt.io Secure Profile shuffled advertisement
type KontaktShuffledAdvertisement struct {
	DeviceModel         uint8
	FirmwareMajor       uint8
	FirmwareMinor       uint8
	BatteryLevel        uint8
	EddystoneNamespace  []byte
	EddystoneInstanceID []byte
}

// KontaktLocationAdvertisement is a structure holding data from Kontakt.io Location advertisement
// Location advertisement is described here: https://developer.kontakt.io/hardware/packets/location/
type KontaktLocationAdvertisement struct {
	TxPower     int8
	BleChannel  uint8
	DeviceModel uint8
	Flags       uint8
	UniqueID    string
}

// TelemetryPID is an identifier of telemetry field.
type TelemetryPID uint16

const (
	// SystemHealth - grouped field https://developer.kontakt.io/hardware/packets/telemetry/#system-health-beacon
	SystemHealth TelemetryPID = 0x01
	// Accelerometer - grouped field https://developer.kontakt.io/hardware/packets/telemetry/#accelerometer
	Accelerometer TelemetryPID = 0x02
	// Sensors - grouped field https://developer.kontakt.io/hardware/packets/telemetry/#accelerometer
	Sensors TelemetryPID = 0x05
	// Acceleration - simple field https://developer.kontakt.io/hardware/packets/telemetry/#acceleration
	Acceleration TelemetryPID = 0x06
	// Movement - simple field https://developer.kontakt.io/hardware/packets/telemetry/#movementfree-fall
	Movement TelemetryPID = 0x07
	// DoubleTap - simple field https://developer.kontakt.io/hardware/packets/telemetry/#double-tap
	DoubleTap TelemetryPID = 0x08
	// LightLevel - simple field https://developer.kontakt.io/hardware/packets/telemetry/#light-level
	LightLevel TelemetryPID = 0x0A
	// Temperature8Bit - simple field https://developer.kontakt.io/hardware/packets/telemetry/#temperature
	Temperature8Bit TelemetryPID = 0x0B
	// Temperature16Bit - simple field https://developer.kontakt.io/hardware/packets/telemetry/#precise-temperature
	Temperature16Bit TelemetryPID = 0x13
	// BatteryLevel - simple field https://developer.kontakt.io/hardware/packets/telemetry/#battery
	BatteryLevel TelemetryPID = 0x0C
	// TimeSinceClick - simple field https://developer.kontakt.io/hardware/packets/telemetry/#button-click
	TimeSinceClick TelemetryPID = 0x0D
	// ClickInfo - simple field https://developer.kontakt.io/hardware/packets/telemetry/#button-click-counter
	ClickInfo TelemetryPID = 0x11
	// UTCTime - simple field https://developer.kontakt.io/hardware/packets/telemetry/#utc-time
	UTCTime TelemetryPID = 0x0F
	// Humidity - simple field https://developer.kontakt.io/hardware/packets/telemetry/#humidity
	Humidity TelemetryPID = 0x12
	// MovementInfo - simple field https://developer.kontakt.io/hardware/packets/telemetry/#movement-counter
	MovementInfo TelemetryPID = 0x16
)

// KontaktTelemetryValue is a container for storing single value of telemetry data
type KontaktTelemetryValue struct {
	PID   uint16
	value []byte
}

// KontaktTelemetryAdvertisement is describing contents of a Kontakt.io telemetry advertisement
type KontaktTelemetryAdvertisement struct {
	Values []KontaktTelemetryValue
}
