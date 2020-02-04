package kontaktparser

const (
	// IBeacon - iBeacon packet
	IBeacon = iota
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
