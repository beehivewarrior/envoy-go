package envoy_go

type SessionResponse struct {
	Message      string `json:"message"`
	SessionId    string `json:"session_id"`
	ManagerToken string `json:"manager_token"`
	Consumer     bool   `json:"is_consumer"`
}

type TokenRequest struct {
	SessionID string `json:"session_id"`
	Serial    string `json:"serial_num"`
	Username  string `json:"username"`
}

type NetworkInterface struct {
	InterfaceType     string `json:"type"`
	InterfaceName     string `json:"interface"`
	MacAddress        string `json:"mac"`
	DHCPServer        bool   `json:"dhcp"`
	IP                string `json:"ip"`
	SignalStrength    int64  `json:"signal_strength"`
	MaxSignalStrength int64  `json:"signal_strength_max"`
	Carrier           bool   `json:"carrier"`
	Supported         bool   `json:"supported"`
	Present           bool   `json:"present"`
	Configured        bool   `json:"configured"`
	Status            string `json:"status"`
}

type Network struct {
	Connected           bool               `json:"web_comm"`
	ReportedToEnlighten bool               `json:"ever_reported_to_enlighten"`
	LastEnlightenReport int64              `json:"last_enlighten_report_time"`
	PrimaryInterface    string             `json:"primary_interface"`
	Interfaces          []NetworkInterface `json:"interfaces"`
}

type WirelessConnection struct {
	SignalStrength    int64  `json:"signal_strength"`
	MaxSignalStrength int64  `json:"signal_strength_max"`
	ConnectionType    string `json:"type"`
	Connected         bool   `json:"connected"`
}

type EnpowerStatus struct {
	Connected  bool   `json:"connected"`
	GridStatus string `json:"grid_status"`
}

type SystemInfo struct {
	BuildEpoch          int64                  `json:"software_build_epoch"`
	Nonvoy              bool                   `json:"is_nonvoy"`
	DBSize              int64                  `json:"db_size"`
	DBUsage             string                 `json:"db_percent_full"`
	Timezone            string                 `json:"timezone"`
	CurrentDate         string                 `json:"current_date"`
	CurrentTime         string                 `json:"current_time"`
	Network             Network                `json:"network"`
	Tariff              string                 `json:"tariff"`
	Comm                map[string]interface{} `json:"-"`
	Alerts              []string               `json:"alerts"`
	Updates             string                 `json:"update_status"`
	WirelessConnections []WirelessConnection   `json:"wireless_connections"`
	Enpower             EnpowerStatus          `json:"enpower"`
}

type Meter struct {
	ID              string   `json:"eid"`
	State           string   `json:"state"`
	Status          string   `json:"meteringStatus"`
	MeasurementType string   `json:"measurementType"`
	Mode            string   `json:"phaseMode"`
	Phases          int      `json:"phaseCount"`
	Flags           []string `json:"statusFlags"`
}

type PowerMeasure struct {
	TimeStamp       int64   `json:"timestamp"`
	EnergyDelivered float64 `json:"actEnergyDlvd"`
	EnergyReceived  float64 `json:"actEnergyRcvd"`
	ApparentEnergy  float64 `json:"apparentEnergy"`
	ReactiveLag     float64 `json:"reactiveEnergyLagg"`
	ReactiveLead    float64 `json:"reactiveEnergyLead"`
	InstantDemand   float64 `json:"instantaneousDemand"`
	ActivePower     float64 `json:"activePower"`
	ApparentPower   float64 `json:"apparentPower"`
	ReactivePower   float64 `json:"reactivePower"`
	PowerFactor     float64 `json:"pwrFactor"`
	Voltage         float64 `json:"voltage"`
	Current         float64 `json:"current"`
	Frequency       float64 `json:"freq"`
}

type Channel struct {
	ChannelID string `json:"eid"`
	TimeStamp int64  `json:"timestamp"`
	PowerMeasure
}

type MeterReading struct {
	MeterID  string    `json:"eid"`
	Channels []Channel `json:"channels"`
	PowerMeasure
}

type InverterReading struct {
	SerialNumber string `json:"SerialNumber"`
	TimeStamp    int64  `json:"timestamp"`
	InverterType string `json:"devType"`
	LastReport   int64  `json:"lastReportWatts"`
	MaxReport    int64  `json:"maxReportWatts"`
}
