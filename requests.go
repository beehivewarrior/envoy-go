package envoy_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func RequestSession(username, password, portal string) (*http.Response, error) {

	client := &http.Client{}

	formFields := url.Values{"user[email]": {username}, "user[password]": {password}}
	resp, err := client.PostForm(portal, formFields)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RequestAuthToken(sessionId, serial, username, portal string) (*http.Response, error) {
	if sessionId == "" {
		return nil, fmt.Errorf("failed to get auth token: session id is empty")
	}

	if serial == "" {
		return nil, fmt.Errorf("failed to get auth token: serial number is empty")
	}

	client := &http.Client{}

	tokenRequest := TokenRequest{SessionID: sessionId, Serial: serial, Username: username}
	tokenRequestJson, err := json.Marshal(tokenRequest)

	if err != nil {
		return nil, err
	}

	resp, err := client.Post(portal, "application/json", bytes.NewBuffer(tokenRequestJson))

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RequestSystemInfo(authToken, gateway string, client *http.Client) (*http.Response, error) {
	if authToken == "" {
		return nil, fmt.Errorf("failed to get system info: auth token is empty - please login")
	}

	if gateway == "" {
		return nil, fmt.Errorf("failed to get system info: gateway is empty")
	}

	if client == nil {
		client = NewGatewayClient()
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/home.json", gateway), nil)

	if err != nil {
		return nil, err
	}

	// Set the auth token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	// Send the request
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RequestMeters(authToken, gateway string, client *http.Client) (*http.Response, error) {
	if authToken == "" {
		return nil, fmt.Errorf("failed to get meters: auth token is empty - please login")
	}

	if gateway == "" {
		return nil, fmt.Errorf("failed to get meters: gateway is empty")
	}

	if client == nil {
		client = NewGatewayClient()
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ivp/meters", gateway), nil)

	if err != nil {
		return nil, err
	}

	// Set the auth token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	// Send the request
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RequestMeterReadings(authToken, gateway string, client *http.Client) (*http.Response, error) {
	if authToken == "" {
		return nil, fmt.Errorf("failed to get meter readings: auth token is empty - please login")
	}

	if gateway == "" {
		return nil, fmt.Errorf("failed to get meter readings: gateway is empty")
	}

	if client == nil {
		client = NewGatewayClient()
	}

	resourcePath := fmt.Sprintf("%s/ivp/meters/readings", gateway)

	resp, err := MakeAuthorizedRequest("GET", resourcePath, authToken, nil, client)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RequestInverterReadings(authToken, gateway string, client *http.Client) (*http.Response, error) {
	if authToken == "" {
		return nil, fmt.Errorf("failed to get inverter reading: auth token is empty - please login")
	}

	if gateway == "" {
		return nil, fmt.Errorf("failed to get inverter reading: gateway is empty")
	}

	if client == nil {
		client = NewGatewayClient()
	}

	resourcePath := fmt.Sprintf("%s/api/v1/production/inverters/summary", gateway)

	resp, err := MakeAuthorizedRequest("GET", resourcePath, authToken, nil, client)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
