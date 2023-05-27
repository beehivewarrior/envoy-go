package envoy_go

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultGateway = "https://envoy.local"
	sessionPortal  = "https://enlighten.enphaseenergy.com/login/login.json?"
	authPortal     = "https://entrez.enphaseenergy.com/tokens"
)

type EnvoyClient struct {
	gateway        string
	sessionGateway string
	authGateway    string
	authToken      string
	username       string
	serial         string
	system         *EnvoySystem
	client         *http.Client
}

func NewEnvoyClient(gateway, sessionURL, authURL, username, serial string) *EnvoyClient {
	if gateway == "" {
		gateway = defaultGateway
	}

	if sessionURL == "" {
		sessionURL = sessionPortal
	}

	if authURL == "" {
		authURL = authPortal
	}

	return &EnvoyClient{
		gateway: gateway, username: username, authToken: "", serial: serial,
		client: NewGatewayClient(), authGateway: authURL, sessionGateway: sessionURL,
	}
}

func (e *EnvoyClient) Gateway() string {
	return e.gateway
}

func (e *EnvoyClient) Client() *http.Client {
	return e.client
}

func (e *EnvoyClient) Close() {
	e.client.CloseIdleConnections()
}

func (e *EnvoyClient) Login(password string) (bool, error) {
	sessionId, err := e.getSessionId(password)

	if err != nil {
		return false, err
	}

	authToken, err := e.getAuthToken(*sessionId)

	if err != nil {
		return false, err
	}

	e.authToken = *authToken

	//// Attempt to get system info to verify login
	//_, err = e.SystemInfo()
	//
	//if err != nil {
	//	return false, err
	//}

	return true, nil
}

func (e *EnvoyClient) Meters() (*[]Meter, error) {

	if e.system != nil && e.system.Meters != nil {
		return e.system.Meters, nil
	}

	meters, err := e.getMeters()

	if err != nil {
		return nil, err
	}

	if e.system == nil {
		e.system = &EnvoySystem{}
	}

	e.system.Meters = meters

	return e.system.Meters, nil
}

func (e *EnvoyClient) ReadInverters() ([]InverterReading, error) {
	return e.getLastInverterReads()
}

func (e *EnvoyClient) ReadMeters() ([]MeterReading, error) {
	return e.getMeterReads()
}

func (e *EnvoyClient) SystemInfo() (*SystemInfo, error) {

	if e.system != nil && e.system.Info != nil {
		return e.system.Info, nil
	}

	systemInfo, err := e.getSystemInfo()

	if err != nil {
		return nil, err
	}

	if e.system == nil {
		e.system = &EnvoySystem{}
	}

	e.system.Info = systemInfo

	return e.system.Info, nil
}

func (e *EnvoyClient) getLastInverterReads() ([]InverterReading, error) {
	resp, err := RequestInverterReadings(e.authToken, e.gateway, e.client)

	if err != nil {
		return nil, err
	}

	// If response is not 200, return error
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get inverter reads: %s", resp.Status)
	}

	var inverterReads []InverterReading

	// Read the response body
	resp, err = LoadResponse(resp, &inverterReads)

	if err != nil {
		return nil, err
	}

	return inverterReads, nil
}

func (e *EnvoyClient) getMeterReads() ([]MeterReading, error) {
	resp, err := RequestMeterReadings(e.authToken, e.gateway, e.client)

	if err != nil {
		return nil, err
	}

	// If response is not 200, return error
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get meter reads: %s", resp.Status)
	}

	var meterReads []MeterReading

	// Read the response body
	resp, err = LoadResponse(resp, &meterReads)

	if err != nil {
		return nil, err
	}

	return meterReads, nil

}

func (e *EnvoyClient) getSessionId(password string) (*string, error) {
	resp, err := RequestSession(e.username, password, e.sessionGateway)

	if err != nil {
		return nil, err
	}

	sessionResponse := &SessionResponse{}

	resp, err = LoadResponse(resp, &sessionResponse)

	if sessionResponse.Message != "success" {
		return nil, fmt.Errorf("failed to get session id: %s", sessionResponse.Message)
	} else if sessionResponse.SessionId == "" {
		return nil, fmt.Errorf("failed to get session id: session id is empty")
	} else {
		return &sessionResponse.SessionId, nil
	}

}

func (e *EnvoyClient) getSystemInfo() (*SystemInfo, error) {
	resp, err := RequestSystemInfo(e.authToken, e.gateway, e.client)

	if err != nil {
		return nil, err
	}

	// If response is not 200, return error
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get system info: %s", resp.Status)
	}

	systemInfo := SystemInfo{}

	// Read the response body
	resp, err = LoadResponse(resp, &systemInfo)

	if err != nil {
		return nil, err
	}

	return &systemInfo, nil
}

func (e *EnvoyClient) getAuthToken(sessionId string) (*string, error) {
	resp, err := RequestAuthToken(sessionId, e.serial, e.username, e.authGateway)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	token := bytes.NewBuffer(body).String()

	e.authToken = token

	return &token, nil

}

func (e *EnvoyClient) getMeters() (*[]Meter, error) {
	resp, err := RequestMeters(e.authToken, e.gateway, e.client)

	if err != nil {
		return nil, err
	}

	// If response is not 200, return error
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get meter info: %s", resp.Status)
	}

	var meters []Meter

	// Read the response body
	resp, err = LoadResponse(resp, &meters)

	if err != nil {
		return nil, err
	}

	return &meters, nil

}
