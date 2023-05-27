package envoy_go

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

func NewGatewayClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Jar: jar, Transport: transport}
}

func MakeAuthorizedRequest(method, url, authToken string, body io.Reader, client *http.Client) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)

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

func ReadResponse(resp *http.Response) ([]byte, error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func LoadResponse(resp *http.Response, responseObj interface{}) (*http.Response, error) {
	body, err := ReadResponse(resp)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &responseObj)

	if err != nil {
		return nil, err
	}

	return resp, nil

}
