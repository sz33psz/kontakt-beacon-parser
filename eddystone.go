package kontaktparser

type EddystoneUIDPacket struct {
	TxPower0M  int8
	Namespace  []byte
	InstanceId []byte
}

type EddystoneURLPacket struct {
	TxPower0M int8
	URL       string
}

type EddystonePlainTLMPacket struct {
	BatteryVoltage     uint16
	Temperature        float64
	AdvertisementCount uint32
	TimeSincePowerOn   float64
}

type EddystoneEncryptedTLMPacket struct {
	Telemetry []byte
	Salt      []byte
	MIC       []byte
}

type EddystoneEIDPacket struct {
	TxPower0M int8
	EID       []byte
}
