package envoy_go_test

import (
	"encoding/json"
	"fmt"
	envoy_go "github.com/beehivewarrior/envoy-go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	muxEnvoy      *http.ServeMux
	envoyServer   *httptest.Server
	muxEnphase    *http.ServeMux
	enphaseServer *httptest.Server

	client *envoy_go.EnvoyClient
)

var systemTestInfo []byte
var loginSuccess []byte
var authorizationToken []byte

func init() {
	if systemInfo, err := os.ReadFile("test/system_info.json"); err == nil {
		systemTestInfo = systemInfo
	} else {
		panic(err)
	}

	if loginSuccessMsg, err := os.ReadFile("test/login_success.json"); err == nil {
		loginSuccess = loginSuccessMsg
	} else {
		panic(err)
	}

	if authToken, err := os.ReadFile("test/authorization_token"); err == nil {
		authorizationToken = authToken
	} else {
		panic(err)
	}
}

func setup(t *testing.T) (func(), error) {
	muxEnphase = http.NewServeMux()
	enphaseServer = httptest.NewServer(muxEnphase)

	muxEnphase.HandleFunc("/login/login.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username := r.FormValue("user[email]")
		password := r.FormValue("user[password]")

		if username == "envoy" && password == "Test1234" {
			http.SetCookie(w, &http.Cookie{Name: "_enlighten_4_session", Value: "1234567890", Path: "/"})
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(loginSuccess)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})

	muxEnphase.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		validSession := false

		tokenRequest := &envoy_go.TokenRequest{}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(body, &tokenRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, c := range r.Cookies() {
			if c.Name == "_enlighten_4_session" && c.Value == "1234567890" {
				validSession = true
			}
		}

		if !validSession {
			if tokenRequest.SessionID != "1234567890" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if serial := tokenRequest.Serial; serial != "9876543210" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if username := tokenRequest.Username; username != "envoy" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			validSession = true
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(authorizationToken)
	})

	muxEnvoy = http.NewServeMux()
	envoyServer = httptest.NewServer(muxEnvoy)

	client = envoy_go.NewEnvoyClient(envoyServer.URL, fmt.Sprintf("%s/login/login.json", enphaseServer.URL),
		fmt.Sprintf("%s/tokens", enphaseServer.URL), "envoy", "9876543210")

	_, err := client.Login("Test1234")

	if !assert.Nil(t, err) {
		return func() {
			enphaseServer.Close()
			envoyServer.Close()
		}, err
	}

	return func() {
		enphaseServer.Close()
		envoyServer.Close()
	}, nil
}

func TestEnvoyClient_Login(t *testing.T) {
	teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	success, err := client.Login("Test1234")

	if !assert.Nil(t, err) || !assert.True(t, success, "Login failed") {
		return
	}
}
